package service

// totp_service.go — TOTP (Time-based One-Time Password) MFA service.
//
// Contains:
//  1. TOTPService  — stateless helper for key generation and validation (unit-testable).
//  2. AuthService methods — EnableTOTP, VerifyTOTP, DisableTOTP (full DB-backed logic).
//  3. EnrollTOTP / ConfirmTOTP — Sprint-1.8 vocabulary aliases.
//  4. AES-256-GCM crypto helpers — totpEncryptionKey, aesGCMEncrypt, aesGCMDecrypt.

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
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ── Standalone TOTPService ────────────────────────────────────────────────────

// TOTPService is a stateless helper for TOTP key generation and validation.
// It does not hold any repository reference; all DB interactions are handled by
// AuthService which delegates to this service for the cryptographic operations.
type TOTPService struct{}

// NewTOTPService returns a new TOTPService instance.
func NewTOTPService() *TOTPService { return &TOTPService{} }

// GenerateKey creates a new TOTP key for the given issuer and account name.
// SecretSize is 32 bytes (256 bits), period 30 s, 6 digits — RFC 6238 defaults.
// Returns the otpauth:// provisioning URI and the raw base32 secret.
func (t *TOTPService) GenerateKey(issuer, accountName string) (provisioningURI, secret string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		SecretSize:  32,
		Period:      30,
		Digits:      6,
	})
	if err != nil {
		logger.Errorf("generate TOTP key: %v", err)
		return "", "", errors.New("generate TOTP key")
	}
	return key.URL(), key.Secret(), nil
}

// Validate checks a 6-digit TOTP code against a plaintext base32 secret.
// Allows ±1 step (±30 s) clock skew tolerance.
func (t *TOTPService) Validate(code, secret string) (bool, error) {
	valid, err := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period: 30,
		Skew:   1,
		Digits: 6,
	})
	if err != nil {
		logger.Errorf("totp validate: %v", err)
		return false, errors.New("totp validate")
	}
	return valid, nil
}

// ── Package-level helpers (used by tests without a TOTPService instance) ─────

// EnrollTOTPForUser generates a TOTP key for the given issuer and account.
// Returns the otpauth:// provisioning URI and raw base32 secret.
func EnrollTOTPForUser(issuer, accountName string) (provisioningURI, secret string, err error) {
	return NewTOTPService().GenerateKey(issuer, accountName)
}

// ValidateTOTPCode validates a TOTP code against a plaintext base32 secret
// with ±1 step (30 s) tolerance.
func ValidateTOTPCode(code, secret string) (bool, error) {
	return NewTOTPService().Validate(code, secret)
}

// ── AuthService TOTP methods ──────────────────────────────────────────────────

// EnableTOTP generates a new TOTP secret for the user, encrypts it with AES-256-GCM,
// stores the ciphertext, and returns the provisioning URI + raw secret for QR code generation.
// totp_enabled stays false until VerifyTOTP confirms the first valid code.
func (s *AuthService) EnableTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("EnableTOTP: user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if user.TotpEnabled {
		return nil, errors.New("TOTP is already enabled for this user")
	}

	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		logger.Errorf("EnableTOTP: generate secret: %v", err)
		return nil, errors.New("generate TOTP secret")
	}
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secretBytes)

	issuer := s.config.JWT.Issuer
	if issuer == "" {
		issuer = "InsureTech"
	}
	accountName := user.MobileNumber
	if user.Email != "" {
		accountName = user.Email
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer: issuer, AccountName: accountName,
		Secret: secretBytes, Period: 30, Digits: 6,
	})
	if err != nil {
		logger.Errorf("EnableTOTP: generate key: %v", err)
		return nil, errors.New("generate TOTP key")
	}

	encSecret, err := aesGCMEncrypt(secret, totpEncryptionKey())
	if err != nil {
		logger.Errorf("EnableTOTP: encrypt secret: %v", err)
		return nil, errors.New("encrypt TOTP secret")
	}
	if err := s.userRepo.UpdateTOTPSecret(ctx, req.UserId, encSecret); err != nil {
		logger.Errorf("EnableTOTP: store secret: %v", err)
		return nil, errors.New("store TOTP secret")
	}
	return &authnservicev1.EnableTOTPResponse{TotpSecret: secret, ProvisioningUri: key.URL()}, nil
}

// VerifyTOTP validates a TOTP code. On first success after EnableTOTP, activates totp_enabled=true.
// If mfa_session_token is provided, consumes it and issues real session tokens (MFA login flow).
func (s *AuthService) VerifyTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("VerifyTOTP: user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if user.TotpSecretEnc == "" {
		return &authnservicev1.VerifyTOTPResponse{Verified: false, Message: "TOTP not configured. Call EnableTOTP first."}, nil
	}

	secret, err := aesGCMDecrypt(user.TotpSecretEnc, totpEncryptionKey())
	if err != nil {
		logger.Errorf("VerifyTOTP: decrypt secret: %v", err)
		return nil, errors.New("decrypt TOTP secret")
	}
	if !totp.Validate(req.TotpCode, secret) {
		return &authnservicev1.VerifyTOTPResponse{Verified: false, Message: "Invalid TOTP code"}, nil
	}

	if !user.TotpEnabled {
		if err := s.userRepo.SetTOTPEnabled(ctx, req.UserId, true); err != nil {
			logger.Errorf("VerifyTOTP: activate: %v", err)
			return nil, errors.New("activate TOTP")
		}
	}

	resp := &authnservicev1.VerifyTOTPResponse{Verified: true, Message: "TOTP verified successfully"}

	if req.MfaSessionToken != "" {
		userID, deviceID, deviceType, ipAddress, consumeErr := s.ConsumeMFASessionToken(ctx, req.MfaSessionToken)
		if consumeErr != nil {
			logger.Errorf("VerifyTOTP: invalid MFA token: %v", consumeErr)
			return nil, errors.New("invalid or expired MFA session token")
		}
		if userID != req.UserId {
			return nil, errors.New("MFA session token user mismatch")
		}
		parsedDeviceType := parseDeviceType(deviceType)
		if parsedDeviceType == authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED {
			return nil, status.Errorf(codes.InvalidArgument, "MFA session token contains invalid device_type %q", deviceType)
		}
		sessionType := mapDeviceTypeToSessionType(parsedDeviceType)
		if sessionType == authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE {
			serverSession, err := s.tokenService.GenerateServerSideSession(ctx, userID, deviceID, parsedDeviceType, ipAddress, "")
			if err != nil {
				logger.Errorf("VerifyTOTP: MFA session: %v", err)
				return nil, errors.New("MFA post-verification session creation failed")
			}
			resp.SessionToken = serverSession.SessionToken
			resp.SessionId = serverSession.SessionID
			resp.CsrfToken = serverSession.CSRFToken
			resp.SessionType = "SERVER_SIDE"
		} else {
			tokens, err := s.tokenService.GenerateJWT(ctx, userID, user.UserType.String(), "", deviceID, parsedDeviceType, ipAddress, "")
			if err != nil {
				logger.Errorf("VerifyTOTP: MFA JWT: %v", err)
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

// DisableTOTP verifies the current code then clears the secret and disables TOTP.
func (s *AuthService) DisableTOTP(ctx context.Context, req *authnservicev1.DisableTOTPRequest) (*authnservicev1.DisableTOTPResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("DisableTOTP: user not found: %v", err)
		return nil, errors.New("user not found")
	}
	if !user.TotpEnabled || user.TotpSecretEnc == "" {
		return nil, errors.New("TOTP is not enabled for this user")
	}
	secret, err := aesGCMDecrypt(user.TotpSecretEnc, totpEncryptionKey())
	if err != nil {
		logger.Errorf("DisableTOTP: decrypt: %v", err)
		return nil, errors.New("decrypt TOTP secret")
	}
	if !totp.Validate(req.TotpCode, secret) {
		return nil, errors.New("invalid TOTP code — cannot disable TOTP")
	}
	if err := s.userRepo.UpdateTOTPSecret(ctx, req.UserId, ""); err != nil {
		logger.Errorf("DisableTOTP: clear secret: %v", err)
		return nil, errors.New("clear TOTP secret")
	}
	if err := s.userRepo.SetTOTPEnabled(ctx, req.UserId, false); err != nil {
		logger.Errorf("DisableTOTP: disable: %v", err)
		return nil, errors.New("disable TOTP")
	}
	return &authnservicev1.DisableTOTPResponse{Message: "TOTP disabled successfully"}, nil
}

// ── AuthService delegation shims (Sprint-1.8 vocabulary) ─────────────────────

// EnrollTOTP is an alias for EnableTOTP using Sprint-1.8 naming convention.
func (s *AuthService) EnrollTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	return s.EnableTOTP(ctx, req)
}

// ConfirmTOTP is an alias for VerifyTOTP using Sprint-1.8 naming convention.
func (s *AuthService) ConfirmTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	return s.VerifyTOTP(ctx, req)
}

// ── AES-256-GCM crypto helpers ────────────────────────────────────────────────

// totpEncryptionKey returns the AES-256 key from env var TOTP_ENCRYPTION_KEY (base64).
// Falls back to a zero key in dev — MUST set TOTP_ENCRYPTION_KEY in production.
func totpEncryptionKey() []byte {
	keyB64 := os.Getenv("TOTP_ENCRYPTION_KEY")
	if keyB64 == "" {
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
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}

// aesGCMDecrypt decrypts base64(nonce+ciphertext) produced by aesGCMEncrypt.
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
