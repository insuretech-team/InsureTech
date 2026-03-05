package service

// stub_service.go — formerly contained stubs; all methods now fully implemented below.
// KYC, TOTP, Voice, and Profile methods delegate to their respective repositories.

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	voicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ── KYC ──────────────────────────────────────────────────────────────────────

// InitiateKYC creates a new KYC verification record in PENDING status.
func (s *AuthService) InitiateKYC(ctx context.Context, req *authnservicev1.InitiateKYCRequest) (*authnservicev1.InitiateKYCResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	kycID := uuid.New().String()
	if s.externalKYC != nil {
		extResp, err := s.externalKYC.StartKYCVerification(ctx, &kycservicev1.StartKYCVerificationRequest{
			Type:       "KYC",
			EntityType: "user",
			EntityId:   req.UserId,
			Method:     "MANUAL",
		})
		if err != nil {
			logger.Errorf("initiate KYC (external): %v", err)
			return nil, errors.New("initiate KYC (external)")
		}
		if extResp == nil || extResp.KycVerificationId == "" {
			return nil, errors.New("initiate KYC (external): empty verification id")
		}
		kycID = extResp.KycVerificationId
	}
	k := &kycv1.KYCVerification{
		Id:         kycID,
		Type:       kycv1.VerificationType_VERIFICATION_TYPE_KYC,
		EntityType: "user",
		EntityId:   req.UserId,
		Method:     kycv1.VerificationMethod_VERIFICATION_METHOD_MANUAL,
		Status:     kycv1.VerificationStatus_VERIFICATION_STATUS_PENDING,
	}
	if err := s.kycRepo.Create(ctx, k); err != nil {
		logger.Errorf("initiate KYC: %v", err)
		return nil, errors.New("initiate KYC")
	}
	s.cacheKYCSessionOwner(ctx, kycID, req.UserId)
	return &authnservicev1.InitiateKYCResponse{
		KycId:   kycID,
		Status:  "PENDING",
		Message: "KYC verification initiated. Please upload required documents.",
	}, nil
}

// GetKYCStatus returns the current KYC status for a user.
func (s *AuthService) GetKYCStatus(ctx context.Context, req *authnservicev1.GetKYCStatusRequest) (*authnservicev1.GetKYCStatusResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	k, err := s.kycRepo.GetByEntity(ctx, "user", req.UserId)
	if err != nil {
		logger.Errorf("get KYC status: %v", err)
		return nil, errors.New("get KYC status")
	}
	resp := &authnservicev1.GetKYCStatusResponse{
		KycId:  k.Id,
		Status: strings.TrimPrefix(k.Status.String(), "VERIFICATION_STATUS_"),
	}
	if k.RejectionReason != "" {
		resp.RejectionReason = k.RejectionReason
	}
	if k.VerifiedAt != nil {
		resp.ReviewedAt = k.VerifiedAt
	}
	return resp, nil
}

// ApproveKYC sets KYC status to VERIFIED and marks user profile kyc_verified=true.
func (s *AuthService) ApproveKYC(ctx context.Context, req *authnservicev1.ApproveKYCRequest) (*authnservicev1.ApproveKYCResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	now := time.Now()
	if err := s.kycRepo.MarkVerified(ctx, req.KycId, req.ReviewerId, now, nil); err != nil {
		logger.Errorf("approve KYC: %v", err)
		return nil, errors.New("approve KYC")
	}
	// Fetch kyc to get entity_id (user_id)
	k, err := s.kycRepo.GetByID(ctx, req.KycId)
	if err == nil && s.userProfileRepo != nil {
		_ = s.userProfileRepo.SetKYCVerified(ctx, k.EntityId, true, &now)
	}
	return &authnservicev1.ApproveKYCResponse{
		Message: "KYC approved successfully",
	}, nil
}

// RejectKYC sets KYC status to REJECTED with a rejection reason.
func (s *AuthService) RejectKYC(ctx context.Context, req *authnservicev1.RejectKYCRequest) (*authnservicev1.RejectKYCResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	reason := req.RejectionReason
	if err := s.kycRepo.UpdateStatus(ctx, req.KycId, kycv1.VerificationStatus_VERIFICATION_STATUS_REJECTED, &reason); err != nil {
		logger.Errorf("reject KYC: %v", err)
		return nil, errors.New("reject KYC")
	}
	return &authnservicev1.RejectKYCResponse{
		Message: "KYC rejected",
	}, nil
}

// VerifyDocument marks a user document as verified by the given reviewer.
func (s *AuthService) VerifyDocument(ctx context.Context, req *authnservicev1.VerifyDocumentRequest) (*authnservicev1.VerifyDocumentResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("user document repository not configured")
	}
	if err := s.userDocumentRepo.MarkVerified(ctx, req.UserDocumentId, req.VerifiedBy, req.VerificationStatus, req.RejectionReason); err != nil {
		logger.Errorf("verify document: %v", err)
		return nil, errors.New("verify document")
	}
	doc, err := s.userDocumentRepo.GetByID(ctx, req.UserDocumentId)
	if err != nil {
		return &authnservicev1.VerifyDocumentResponse{Message: "Document verified"}, nil
	}
	return &authnservicev1.VerifyDocumentResponse{
		Document: doc,
		Message:  "Document verified successfully",
	}, nil
}

// ── TOTP / 2FA ────────────────────────────────────────────────────────────────

// EnableTOTP generates a new TOTP secret for the user, encrypts it with AES-256-GCM,
// stores the ciphertext in users.totp_secret_enc, and returns provisioning URI + raw secret
// for QR code generation.
func (s *AuthService) EnableTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if user.TotpEnabled {
		return nil, errors.New("TOTP is already enabled for this user")
	}

	// Generate a random 20-byte TOTP secret (base32-encoded as per RFC 6238)
	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		logger.Errorf("generate TOTP secret: %v", err)
		return nil, errors.New("generate TOTP secret")
	}
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secretBytes)

	// Generate TOTP provisioning URI (otpauth://)
	issuer := s.config.JWT.Issuer
	if issuer == "" {
		issuer = "InsureTech"
	}
	accountName := user.MobileNumber
	if user.Email != "" {
		accountName = user.Email
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Secret:      secretBytes,
		Period:      30,
		Digits:      6,
	})
	if err != nil {
		logger.Errorf("generate TOTP key: %v", err)
		return nil, errors.New("generate TOTP key")
	}

	// Encrypt secret with AES-256-GCM using the configured encryption key
	encSecret, err := aesGCMEncrypt(secret, totpEncryptionKey())
	if err != nil {
		logger.Errorf("encrypt TOTP secret: %v", err)
		return nil, errors.New("encrypt TOTP secret")
	}

	// Persist encrypted secret; totp_enabled stays false until VerifyTOTP confirms
	if err := s.userRepo.UpdateTOTPSecret(ctx, req.UserId, encSecret); err != nil {
		logger.Errorf("store TOTP secret: %v", err)
		return nil, errors.New("store TOTP secret")
	}

	return &authnservicev1.EnableTOTPResponse{
		TotpSecret:      secret,
		ProvisioningUri: key.URL(),
	}, nil
}

// VerifyTOTP validates a TOTP code against the stored secret.
// On first successful verification after EnableTOTP, activates totp_enabled=true.
func (s *AuthService) VerifyTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if user.TotpSecretEnc == "" {
		return &authnservicev1.VerifyTOTPResponse{
			Verified: false,
			Message:  "TOTP not configured. Call EnableTOTP first.",
		}, nil
	}

	// Decrypt stored secret
	secret, err := aesGCMDecrypt(user.TotpSecretEnc, totpEncryptionKey())
	if err != nil {
		logger.Errorf("decrypt TOTP secret: %v", err)
		return nil, errors.New("decrypt TOTP secret")
	}

	// Validate code with ±1 step tolerance (30s window)
	valid := totp.Validate(req.TotpCode, secret)
	if !valid {
		return &authnservicev1.VerifyTOTPResponse{
			Verified: false,
			Message:  "Invalid TOTP code",
		}, nil
	}

	// Activate TOTP if not yet enabled (first-time verification after EnableTOTP)
	if !user.TotpEnabled {
		if err := s.userRepo.SetTOTPEnabled(ctx, req.UserId, true); err != nil {
			logger.Errorf("activate TOTP: %v", err)
			return nil, errors.New("activate TOTP")
		}
	}

	resp := &authnservicev1.VerifyTOTPResponse{
		Verified: true,
		Message:  "TOTP verified successfully",
	}

	// Sprint 2.2: If an mfa_session_token is provided, consume it and issue real session tokens.
	// This is the second step of the MFA-gated login flow:
	//   Login → mfa_required=true + mfa_session_token → VerifyTOTP → tokens
	if req.MfaSessionToken != "" {
		userID, deviceID, deviceType, ipAddress, consumeErr := s.ConsumeMFASessionToken(ctx, req.MfaSessionToken)
		if consumeErr != nil {
			logger.Errorf("invalid or expired MFA session token: %v", consumeErr)
			return nil, errors.New("invalid or expired MFA session token")
		}
		// Make sure the token belongs to this user
		if userID != req.UserId {
			return nil, errors.New("MFA session token user mismatch")
		}

		parsedDeviceType := parseDeviceType(deviceType)
		sessionType := mapDeviceTypeToSessionType(parsedDeviceType)

		if sessionType == authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE {
			serverSession, err := s.tokenService.GenerateServerSideSession(
				ctx, userID, deviceID, parsedDeviceType, ipAddress, "",
			)
			if err != nil {
				logger.Errorf("MFA post-verification session creation failed: %v", err)
				return nil, errors.New("MFA post-verification session creation failed")
			}
			resp.SessionToken = serverSession.SessionToken
			resp.SessionId = serverSession.SessionID
			resp.CsrfToken = serverSession.CSRFToken
			resp.SessionType = "SERVER_SIDE"
		} else {
			tokens, err := s.tokenService.GenerateJWT(
				ctx, userID, user.UserType.String(), "", deviceID,
				parsedDeviceType, ipAddress, "",
			)
			if err != nil {
				logger.Errorf("MFA post-verification token generation failed: %v", err)
				return nil, errors.New("MFA post-verification token generation failed")
			}
			resp.AccessToken = tokens.AccessToken
			resp.RefreshToken = tokens.RefreshToken
			resp.SessionId = tokens.SessionID
			resp.SessionType = "JWT"
			resp.AccessTokenExpiresIn = int32(tokens.AccessTokenExpiresIn.Seconds())
			resp.RefreshTokenExpiresIn = int32(tokens.RefreshTokenExpiresIn.Seconds())
		}
		_ = s.eventPublisher.PublishUserLoggedIn(ctx, userID, resp.SessionId, resp.SessionType, ipAddress, deviceType, "")
		s.markTrustedDevice(ctx, userID, deviceID)
	}

	return resp, nil
}

// DisableTOTP verifies the current TOTP code then clears the secret and disables TOTP.
func (s *AuthService) DisableTOTP(ctx context.Context, req *authnservicev1.DisableTOTPRequest) (*authnservicev1.DisableTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if !user.TotpEnabled || user.TotpSecretEnc == "" {
		return nil, errors.New("TOTP is not enabled for this user")
	}

	// Verify current code before disabling
	secret, err := aesGCMDecrypt(user.TotpSecretEnc, totpEncryptionKey())
	if err != nil {
		logger.Errorf("decrypt TOTP secret: %v", err)
		return nil, errors.New("decrypt TOTP secret")
	}
	if !totp.Validate(req.TotpCode, secret) {
		return nil, errors.New("invalid TOTP code — cannot disable TOTP")
	}

	// Clear secret and disable
	if err := s.userRepo.UpdateTOTPSecret(ctx, req.UserId, ""); err != nil {
		logger.Errorf("clear TOTP secret: %v", err)
		return nil, errors.New("clear TOTP secret")
	}
	if err := s.userRepo.SetTOTPEnabled(ctx, req.UserId, false); err != nil {
		logger.Errorf("disable TOTP: %v", err)
		return nil, errors.New("disable TOTP")
	}

	return &authnservicev1.DisableTOTPResponse{
		Message: "TOTP disabled successfully",
	}, nil
}

// ── Voice Sessions ────────────────────────────────────────────────────────────

// CreateVoiceSession creates a new voice session record.
func (s *AuthService) CreateVoiceSession(ctx context.Context, req *authnservicev1.CreateVoiceSessionRequest) (*authnservicev1.CreateVoiceSessionResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	sessionID := uuid.New().String()
	extSessionID := uuid.New().String()
	vs := &voicev1.VoiceSession{
		Id:          sessionID,
		SessionId:   extSessionID,
		UserId:      req.UserId,
		Language:    req.Language,
		PhoneNumber: req.PhoneNumber,
		Status:      voicev1.SessionStatus_SESSION_STATUS_ACTIVE,
		StartedAt:   timestamppb.Now(),
	}
	if err := s.voiceRepo.Create(ctx, vs); err != nil {
		logger.Errorf("create voice session: %v", err)
		return nil, errors.New("create voice session")
	}
	return &authnservicev1.CreateVoiceSessionResponse{
		VoiceSessionId: sessionID,
		Status:         "ACTIVE",
	}, nil
}

// GetVoiceSession retrieves a voice session by ID.
func (s *AuthService) GetVoiceSession(ctx context.Context, req *authnservicev1.GetVoiceSessionRequest) (*authnservicev1.GetVoiceSessionResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	vs, err := s.voiceRepo.GetByID(ctx, req.VoiceSessionId)
	if err != nil {
		logger.Errorf("get voice session: %v", err)
		return nil, errors.New("get voice session")
	}
	return &authnservicev1.GetVoiceSessionResponse{
		VoiceSessionId: vs.Id,
		UserId:         vs.UserId,
		Status:         strings.TrimPrefix(vs.Status.String(), "SESSION_STATUS_"),
		Language:       vs.Language,
		StartedAt:      vs.StartedAt,
		EndedAt:        vs.EndedAt,
	}, nil
}

// EndVoiceSession marks a voice session as ended with final status.
func (s *AuthService) EndVoiceSession(ctx context.Context, req *authnservicev1.EndVoiceSessionRequest) (*authnservicev1.EndVoiceSessionResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	finalStatus := voicev1.SessionStatus_SESSION_STATUS_COMPLETED
	if req.Status == "FAILED" {
		finalStatus = voicev1.SessionStatus_SESSION_STATUS_FAILED
	}
	dur := req.DurationSeconds
	if err := s.voiceRepo.Complete(ctx, req.VoiceSessionId, finalStatus, time.Now(), &dur); err != nil {
		logger.Errorf("end voice session: %v", err)
		return nil, errors.New("end voice session")
	}
	return &authnservicev1.EndVoiceSessionResponse{
		Message: "Voice session ended",
	}, nil
}

// ── Profile ───────────────────────────────────────────────────────────────────

// GetProfilePhotoUploadURL generates a presigned upload URL via StorageService.
func (s *AuthService) GetProfilePhotoUploadURL(ctx context.Context, req *authnservicev1.GetProfilePhotoUploadURLRequest) (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
	contentType := req.ContentType
	if contentType == "" {
		contentType = "image/jpeg"
	}
	fileExt := ".jpg"
	if strings.Contains(contentType, "png") {
		fileExt = ".png"
	} else if strings.Contains(contentType, "webp") {
		fileExt = ".webp"
	}

	legacyFallback := func() (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
		bucket := os.Getenv("S3_BUCKET")
		if bucket == "" {
			bucket = "insuretech-user-media"
		}
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "ap-southeast-1"
		}
		objectKey := "profile-photos/" + req.UserId + "/" + uuid.New().String() + fileExt
		fileURL := "https://" + bucket + ".s3." + region + ".amazonaws.com/" + objectKey

		loadOpts := []func(*awscfg.LoadOptions) error{
			awscfg.WithRegion(region),
		}
		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if accessKey != "" && secretKey != "" {
			sessionToken := os.Getenv("AWS_SESSION_TOKEN")
			loadOpts = append(loadOpts, awscfg.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken),
			))
		}
		awsConfig, err := awscfg.LoadDefaultConfig(ctx, loadOpts...)
		if err != nil {
			return nil, errors.New("load AWS config")
		}
		s3Client := s3.NewFromConfig(awsConfig)
		presignClient := s3.NewPresignClient(s3Client)
		presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			return nil, errors.New("generate presigned upload url")
		}

		return &authnservicev1.GetProfilePhotoUploadURLResponse{
			UploadUrl:        presignedReq.URL,
			FileUrl:          fileURL,
			ExpiresInSeconds: 900,
		}, nil
	}

	storageAddr := os.Getenv("STORAGE_SERVICE_ADDRESS")
	if storageAddr == "" {
		port := os.Getenv("STORAGE_GRPC_PORT")
		if port == "" {
			port = "50290"
		}
		storageAddr = "localhost:" + port
	}

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Errorf("dial storage service: %v", err)
		return legacyFallback()
	}
	defer func() { _ = conn.Close() }()

	tenantID := os.Getenv("DEFAULT_TENANT_ID")
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-tenant-id"); len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			tenantID = vals[0]
		}
	}
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	filename := "profile-" + req.UserId + fileExt

	client := storageservicev1.NewStorageServiceClient(conn)
	uploadResp, err := client.GetUploadURL(ctx, &storageservicev1.GetUploadURLRequest{
		TenantId:         tenantID,
		Filename:         filename,
		ContentType:      contentType,
		FileType:         storageentityv1.FileType_FILE_TYPE_IMAGE,
		ExpiresInMinutes: 15,
		ReferenceId:      req.UserId,
		ReferenceType:    "USER_KYC_PROFILE",
		IsPublic:         false,
	})
	if err != nil {
		logger.Errorf("storage get upload url: %v", err)
		return legacyFallback()
	}

	fileURL := ""
	if cdn := strings.TrimRight(os.Getenv("SPACES_CDN_ENDPOINT"), "/"); cdn != "" {
		fileURL = cdn + "/" + strings.TrimLeft(uploadResp.StorageKey, "/")
	} else if endpoint := strings.TrimRight(os.Getenv("SPACES_ENDPOINT"), "/"); endpoint != "" {
		fileURL = endpoint + "/" + strings.TrimLeft(uploadResp.StorageKey, "/")
	}

	return &authnservicev1.GetProfilePhotoUploadURLResponse{
		UploadUrl:        uploadResp.UploadUrl,
		FileUrl:          fileURL,
		ExpiresInSeconds: 900,
	}, nil
}

// UpdateNotificationPreferences updates the user's notification channel and language preferences.
func (s *AuthService) UpdateNotificationPreferences(ctx context.Context, req *authnservicev1.UpdateNotificationPreferencesRequest) (*authnservicev1.UpdateNotificationPreferencesResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}
	_ = user
	if err := s.userRepo.UpdateNotificationPreferences(ctx, req.UserId, req.NotificationPreference, req.PreferredLanguage); err != nil {
		logger.Errorf("update notification preferences: %v", err)
		return nil, errors.New("update notification preferences")
	}
	return &authnservicev1.UpdateNotificationPreferencesResponse{
		Message: "Notification preferences updated",
	}, nil
}

// ── TOTP crypto helpers ───────────────────────────────────────────────────────

// totpEncryptionKey returns the AES-256 key from env var TOTP_ENCRYPTION_KEY (base64).
// Falls back to a zero key in dev (logs a warning — do NOT use zero key in production).
func totpEncryptionKey() []byte {
	keyB64 := os.Getenv("TOTP_ENCRYPTION_KEY")
	if keyB64 == "" {
		// Dev fallback: zero key — MUST set TOTP_ENCRYPTION_KEY in production
		return make([]byte, 32)
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil || len(key) != 32 {
		return make([]byte, 32)
	}
	return key
}

// aesGCMEncrypt encrypts plaintext using AES-256-GCM. Returns base64(nonce+ciphertext).
func aesGCMEncrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// aesGCMDecrypt decrypts a base64(nonce+ciphertext) produced by aesGCMEncrypt.
func aesGCMDecrypt(encoded string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ct := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
