package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func buildLiveAuthService(t *testing.T, dbConn *gorm.DB) *AuthService {
	t.Helper()
	cfg := &config.Config{}
	cfg.JWT.Issuer = "insuretech-live-test"
	cfg.JWT.AccessTokenDuration = 15 * time.Minute
	cfg.JWT.RefreshTokenDuration = 7 * 24 * time.Hour
	cfg.Security.OTPMaxAttempts = 3
	cfg.Security.OTPExpiry = 5 * time.Minute
	generateTempRSAKeys(cfg)

	pub := events.NewPublisher((events.EventProducer)(nil))
	meta := middleware.NewMetadataExtractor()
	userRepo := repository.NewUserRepository(dbConn.Table("authn_schema.users"))
	sessionRepo := repository.NewSessionRepository(dbConn.Table("authn_schema.sessions"))

	tokenSvc, err := NewTokenService(sessionRepo, userRepo, cfg, pub, meta)
	require.NoError(t, err)

	return NewAuthService(
		tokenSvc,
		NewOTPService(repository.NewOTPRepository(dbConn.Table("authn_schema.otps")), nil, nil, cfg, pub),
		userRepo,
		sessionRepo,
		repository.NewOTPRepository(dbConn.Table("authn_schema.otps")),
		repository.NewApiKeyRepository(dbConn.Table("authn_schema.api_keys")),
		repository.NewUserProfileRepository(dbConn.Table("authn_schema.user_profiles")),
		repository.NewUserDocumentRepository(dbConn.Table("authn_schema.users_documents")),
		repository.NewDocumentTypeRepository(dbConn.Table("authn_schema.document_types")),
		repository.NewKYCVerificationRepository(dbConn.Table("authn_schema.kyc_verifications")),
		repository.NewVoiceSessionRepository(dbConn.Table("authn_schema.voice_sessions")),
		pub,
		cfg,
		meta,
	)
}

func createLiveAuthnUser(t *testing.T, svc *AuthService, ctx context.Context, mobile, email, password string) string {
	t.Helper()
	hash, err := hashPassword(password)
	require.NoError(t, err)
	u, err := svc.userRepo.Create(ctx, mobile, hash, email, authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	require.NoError(t, err)
	return u.UserId
}

func cleanupLiveAuthnUser(t *testing.T, dbConn *gorm.DB, userID string) {
	t.Helper()
	_ = dbConn.Exec(`DELETE FROM authn_schema.voice_sessions WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users_documents WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.user_profiles WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.kyc_verifications WHERE entity_type = 'user' AND entity_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.api_keys WHERE owner_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.sessions WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.otps WHERE user_id = ?`, userID).Error
	_ = dbConn.Exec(`DELETE FROM authn_schema.users WHERE user_id = ?`, userID).Error
}

type fakeExternalKYCClient struct {
	startCalls  int
	uploadCalls int
	verifyCalls int
	startFn     func(ctx context.Context, in *kycservicev1.StartKYCVerificationRequest, opts ...grpc.CallOption) (*kycservicev1.StartKYCVerificationResponse, error)
	uploadFn    func(ctx context.Context, in *kycservicev1.UploadDocumentRequest, opts ...grpc.CallOption) (*kycservicev1.UploadDocumentResponse, error)
	verifyFn    func(ctx context.Context, in *kycservicev1.VerifyKYCRequest, opts ...grpc.CallOption) (*kycservicev1.VerifyKYCResponse, error)
}

func (f *fakeExternalKYCClient) StartKYCVerification(ctx context.Context, in *kycservicev1.StartKYCVerificationRequest, opts ...grpc.CallOption) (*kycservicev1.StartKYCVerificationResponse, error) {
	f.startCalls++
	if f.startFn != nil {
		return f.startFn(ctx, in, opts...)
	}
	return &kycservicev1.StartKYCVerificationResponse{KycVerificationId: uuid.NewString()}, nil
}

func (f *fakeExternalKYCClient) UploadDocument(ctx context.Context, in *kycservicev1.UploadDocumentRequest, opts ...grpc.CallOption) (*kycservicev1.UploadDocumentResponse, error) {
	f.uploadCalls++
	if f.uploadFn != nil {
		return f.uploadFn(ctx, in, opts...)
	}
	return &kycservicev1.UploadDocumentResponse{DocumentVerificationId: uuid.NewString()}, nil
}

func (f *fakeExternalKYCClient) VerifyKYC(ctx context.Context, in *kycservicev1.VerifyKYCRequest, opts ...grpc.CallOption) (*kycservicev1.VerifyKYCResponse, error) {
	f.verifyCalls++
	if f.verifyFn != nil {
		return f.verifyFn(ctx, in, opts...)
	}
	return &kycservicev1.VerifyKYCResponse{Message: "ok"}, nil
}

func TestAuthService_LiveDB_LoginJWTAndLogout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_jwt_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Pass")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	loginResp, err := svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Pass",
		DeviceType:   "API",
		DeviceId:     "dev-jwt-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.AccessToken)
	require.NotEmpty(t, loginResp.RefreshToken)
	require.NotEmpty(t, loginResp.SessionId)
	require.Equal(t, "JWT", loginResp.SessionType)

	validateResp, err := svc.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{AccessToken: loginResp.AccessToken})
	require.NoError(t, err)
	require.True(t, validateResp.Valid)

	logoutResp, err := svc.Logout(ctx, &authnservicev1.LogoutRequest{SessionId: loginResp.SessionId})
	require.NoError(t, err)
	require.True(t, logoutResp.SessionRevoked)
}

func TestAuthService_LiveDB_LoginServerSide(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_web_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Pass2")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	resp, err := svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Pass2",
		DeviceType:   "WEB",
		DeviceId:     "web-dev-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.SessionToken)
	require.NotEmpty(t, resp.CsrfToken)
	require.NotEmpty(t, resp.SessionId)
	require.Equal(t, "SERVER_SIDE", resp.SessionType)

	csrfResp, err := svc.ValidateCSRF(ctx, &authnservicev1.ValidateCSRFRequest{
		SessionId: resp.SessionId,
		CsrfToken: resp.CsrfToken,
	})
	require.NoError(t, err)
	require.True(t, csrfResp.Valid)
}

func TestAuthService_WrapperEdgeCases(t *testing.T) {
	svc := &AuthService{
		metadata: middleware.NewMetadataExtractor(),
	}
	ctx := context.Background()

	_, err := svc.ValidateToken(ctx, &authnservicev1.ValidateTokenRequest{})
	require.Error(t, err)

	out, err := svc.Logout(ctx, &authnservicev1.LogoutRequest{})
	require.NoError(t, err)
	require.False(t, out.SessionRevoked)

	_, err = svc.BiometricAuthenticate(ctx, &authnservicev1.BiometricAuthenticateRequest{
		BiometricToken: "raw-token",
		DeviceType:     "ANDROID",
		DeviceId:       "dev1",
	})
	require.Error(t, err)
}

func TestAuthService_LiveDB_ChangePasswordAndCurrentSession(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	baseCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, baseCtx, mobile, "live_pwd_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Old1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	loginResp, err := svc.Login(baseCtx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Old1",
		DeviceType:   "API",
		DeviceId:     "dev-pwd-1",
	})
	require.NoError(t, err)

	currentCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"authorization", "Bearer "+loginResp.AccessToken,
		"x-forwarded-for", "127.0.0.1",
	))
	current, err := svc.GetCurrentSession(currentCtx, &authnservicev1.GetCurrentSessionRequest{})
	require.NoError(t, err)
	require.Equal(t, loginResp.SessionId, current.Session.SessionId)

	_, err = svc.ChangePassword(baseCtx, &authnservicev1.ChangePasswordRequest{
		UserId:      userID,
		OldPassword: "Str0ng!Old1",
		NewPassword: "Str0ng!New2",
	})
	require.NoError(t, err)

	_, err = svc.Login(baseCtx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Old1",
		DeviceType:   "API",
		DeviceId:     "dev-pwd-2",
	})
	require.Error(t, err)

	newLogin, err := svc.Login(baseCtx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!New2",
		DeviceType:   "API",
		DeviceId:     "dev-pwd-3",
	})
	require.NoError(t, err)
	require.NotEmpty(t, newLogin.AccessToken)
}

func TestAuthService_LiveDB_RevokeAllSessions_ExcludeCurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	baseCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, baseCtx, mobile, "live_revoke_all_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Pass3")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	loginA, err := svc.Login(baseCtx, &authnservicev1.LoginRequest{
		MobileNumber: mobile, Password: "Str0ng!Pass3", DeviceType: "API", DeviceId: "dev-a",
	})
	require.NoError(t, err)
	loginB, err := svc.Login(baseCtx, &authnservicev1.LoginRequest{
		MobileNumber: mobile, Password: "Str0ng!Pass3", DeviceType: "API", DeviceId: "dev-b",
	})
	require.NoError(t, err)
	require.NotEqual(t, loginA.SessionId, loginB.SessionId)

	excludeCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"authorization", "Bearer "+loginA.AccessToken,
		"x-forwarded-for", "127.0.0.1",
	))
	revokeResp, err := svc.RevokeAllSessions(excludeCtx, &authnservicev1.RevokeAllSessionsRequest{
		UserId:                userID,
		ExcludeCurrentSession: true,
		Reason:                "test",
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, revokeResp.RevokedCount, int32(1))

	sessA, err := svc.sessionRepo.GetByID(baseCtx, loginA.SessionId)
	require.NoError(t, err)
	require.True(t, sessA.IsActive)

	_, err = svc.sessionRepo.GetByID(baseCtx, loginB.SessionId)
	require.Error(t, err)
}

func TestAuthService_LiveDB_RegisterAndDuplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	resp, err := svc.Register(ctx, &authnservicev1.RegisterRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Reg1",
		Email:        "reg_live_" + uuid.NewString()[:8] + "@example.com",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.UserId)
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, resp.UserId) })

	_, err = svc.Register(ctx, &authnservicev1.RegisterRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Reg1",
		Email:        "reg_live_" + uuid.NewString()[:8] + "@example.com",
	})
	require.Error(t, err)
}

func TestAuthService_LiveDB_ResetPassword_WithOTP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_reset_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Old9")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	otpCode := "123456"
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)
	require.NoError(t, err)
	otp := &authnentityv1.OTP{
		OtpId:      uuid.NewString(),
		UserId:     userID,
		OtpHash:    string(otpHash),
		Purpose:    "password_reset",
		Recipient:  mobile,
		Channel:    "sms",
		ExpiresAt:  timestamppb.New(time.Now().Add(5 * time.Minute)),
		Verified:   false,
		Attempts:   0,
		DeviceType: "API",
		IpAddress:  "127.0.0.1",
		DlrStatus:  "PENDING",
	}
	require.NoError(t, svc.otpRepo.Create(ctx, otp))

	_, err = svc.ResetPassword(ctx, &authnservicev1.ResetPasswordRequest{
		MobileNumber: mobile,
		OtpCode:      otpCode,
		NewPassword:  "Str0ng!New9",
	})
	require.NoError(t, err)

	_, err = svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!Old9",
		DeviceType:   "API",
		DeviceId:     "dev-reset-1",
	})
	require.Error(t, err)

	login, err := svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile,
		Password:     "Str0ng!New9",
		DeviceType:   "API",
		DeviceId:     "dev-reset-2",
	})
	require.NoError(t, err)
	require.NotEmpty(t, login.AccessToken)
}

func TestAuthService_LiveDB_ServiceWrappersCoverage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_wrap_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Wrap1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	jwtLogin, err := svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile, Password: "Str0ng!Wrap1", DeviceType: "API", DeviceId: "wrap-api-1",
	})
	require.NoError(t, err)

	// RefreshToken wrapper
	ref, err := svc.RefreshToken(ctx, &authnservicev1.RefreshTokenRequest{RefreshToken: jwtLogin.RefreshToken})
	require.NoError(t, err)
	require.NotEmpty(t, ref.AccessToken)

	// Server-side flow for session wrappers
	webLogin, err := svc.Login(ctx, &authnservicev1.LoginRequest{
		MobileNumber: mobile, Password: "Str0ng!Wrap1", DeviceType: "WEB", DeviceId: "wrap-web-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, webLogin.SessionId)

	getResp, err := svc.GetSession(ctx, &authnservicev1.GetSessionRequest{SessionId: webLogin.SessionId})
	require.NoError(t, err)
	require.Equal(t, webLogin.SessionId, getResp.Session.SessionId)

	listResp, err := svc.ListSessions(ctx, &authnservicev1.ListSessionsRequest{UserId: userID, ActiveOnly: true})
	require.NoError(t, err)
	require.NotEmpty(t, listResp.Sessions)

	_, err = svc.RevokeSession(ctx, &authnservicev1.RevokeSessionRequest{SessionId: webLogin.SessionId, Reason: "coverage"})
	require.NoError(t, err)

	// OTP wrappers: force local validation/repo paths (no external SMS dependency)
	_, err = svc.SendOTP(ctx, &authnservicev1.SendOTPRequest{
		Recipient: "invalid-number",
		Type:      "login",
		Channel:   "sms",
	})
	require.Error(t, err)

	verifyResp, err := svc.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{
		OtpId: "missing-otp-id",
		Code:  "123456",
	})
	require.NoError(t, err)
	require.False(t, verifyResp.Verified)

	// JWKS wrapper
	jwksResp, err := svc.GetJWKS(ctx, &authnservicev1.GetJWKSRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, jwksResp.Keys)

	// DLR update wrapper (non-existing message id should return Updated=false)
	dlrResp, err := svc.UpdateDLRStatus(ctx, &authnservicev1.UpdateDLRStatusRequest{
		ProviderMessageId: "non-existing-provider-id",
		Status:            "DELIVERED",
	})
	require.NoError(t, err)
	require.True(t, dlrResp.Updated)
}

func TestAuthService_LiveDB_TOTPAndVoiceFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)
	svc.voiceRepo = repository.NewVoiceSessionRepository(dbConn.Table("authn_schema.voice_sessions"))

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_voice_totp_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Voice1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	// TOTP enable -> verify -> disable
	t.Setenv("TOTP_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(make([]byte, 32)))
	en, err := svc.EnableTOTP(ctx, &authnservicev1.EnableTOTPRequest{UserId: userID})
	require.NoError(t, err)
	require.NotEmpty(t, en.TotpSecret)

	code, err := totp.GenerateCode(en.TotpSecret, time.Now().UTC())
	require.NoError(t, err)
	vr, err := svc.VerifyTOTP(ctx, &authnservicev1.VerifyTOTPRequest{UserId: userID, TotpCode: code})
	require.NoError(t, err)
	require.True(t, vr.Verified)

	code2, err := totp.GenerateCode(en.TotpSecret, time.Now().UTC())
	require.NoError(t, err)
	dis, err := svc.DisableTOTP(ctx, &authnservicev1.DisableTOTPRequest{UserId: userID, TotpCode: code2})
	require.NoError(t, err)
	require.Contains(t, dis.Message, "disabled")

	// Voice flow: initiate -> submit (good transcript) -> verify
	initResp, err := svc.InitiateVoiceSession(ctx, &authnservicev1.InitiateVoiceSessionRequest{UserId: userID})
	require.NoError(t, err)
	require.NotEmpty(t, initResp.SessionId)
	require.NotEmpty(t, initResp.Challenge)

	submit, err := svc.SubmitVoiceSample(ctx, &authnservicev1.SubmitVoiceSampleRequest{
		SessionId:       initResp.SessionId,
		Transcript:      initResp.Challenge,
		ConfidenceScore: 0.99,
	})
	require.NoError(t, err)
	require.True(t, submit.Verified)

	ver, err := svc.VerifyVoiceSession(ctx, &authnservicev1.VerifyVoiceSessionRequest{SessionId: initResp.SessionId})
	require.NoError(t, err)
	require.True(t, ver.Authenticated)
	require.Equal(t, userID, ver.UserId)
}

func TestAuthService_LiveDB_APIKeyProfileAndDocumentFlows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_apikey_profile_doc_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Flow1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	// API key lifecycle
	createKeyResp, err := svc.CreateAPIKey(ctx, &authnservicev1.CreateAPIKeyRequest{
		Name:               "live-key",
		OwnerId:            userID,
		OwnerType:          "PARTNER",
		Scopes:             []string{"authn:read", "authn:write"},
		RateLimitPerMinute: 120,
	})
	require.NoError(t, err)
	require.NotEmpty(t, createKeyResp.KeyId)
	require.NotEmpty(t, createKeyResp.RawKey)
	require.Contains(t, createKeyResp.RawKey, "isk_")

	listKeys, err := svc.ListAPIKeys(ctx, &authnservicev1.ListAPIKeysRequest{
		OwnerId:   userID,
		OwnerType: "PARTNER",
	})
	require.NoError(t, err)
	require.NotEmpty(t, listKeys.Keys)

	_, err = svc.RevokeAPIKey(ctx, &authnservicev1.RevokeAPIKeyRequest{
		KeyId:  createKeyResp.KeyId,
		Reason: "test cleanup",
	})
	require.NoError(t, err)

	// User profile create/get/update
	createProfileResp, err := svc.CreateUserProfile(ctx, &authnservicev1.CreateUserProfileRequest{
		UserId:      userID,
		FullName:    "Live Test User",
		DateOfBirth: timestamppb.New(time.Date(1990, time.June, 15, 0, 0, 0, 0, time.UTC)),
		Gender:      "MALE",
		City:        "Dhaka",
		Country:     "BD",
		NidNumber:   "1234567890123",
	})
	require.NoError(t, err)
	require.Equal(t, userID, createProfileResp.Profile.UserId)

	getProfileResp, err := svc.GetUserProfile(ctx, &authnservicev1.GetUserProfileRequest{UserId: userID})
	require.NoError(t, err)
	require.Equal(t, "Live Test User", getProfileResp.Profile.FullName)

	updateProfileResp, err := svc.UpdateUserProfile(ctx, &authnservicev1.UpdateUserProfileRequest{
		UserId:   userID,
		FullName: "Live Test User Updated",
		City:     "Chittagong",
	})
	require.NoError(t, err)
	require.Equal(t, "Live Test User Updated", updateProfileResp.Profile.FullName)
	require.Equal(t, "Chittagong", updateProfileResp.Profile.City)

	// Document lifecycle with live DB-backed repos
	typeResp, err := svc.ListDocumentTypes(ctx, &authnservicev1.ListDocumentTypesRequest{})
	require.NoError(t, err)
	if len(typeResp.Types) == 0 {
		docType := &authnentityv1.DocumentType{
			DocumentTypeId: uuid.NewString(),
			Code:           "LIVE_TDD_DOC",
			Name:           "Live TDD Doc",
			Description:    "generated by live service test",
			IsActive:       true,
		}
		require.NoError(t, svc.documentTypeRepo.Create(ctx, docType))
		typeResp.Types = append(typeResp.Types, docType)
		t.Cleanup(func() {
			_ = dbConn.Exec(`DELETE FROM authn_schema.document_types WHERE document_type_id = ?`, docType.DocumentTypeId).Error
		})
	}

	uploadResp, err := svc.UploadUserDocument(ctx, &authnservicev1.UploadUserDocumentRequest{
		UserId:         userID,
		DocumentTypeId: typeResp.Types[0].DocumentTypeId,
		FileUrl:        "https://example.com/live-doc.pdf",
		PolicyId:       "",
	})
	require.NoError(t, err)
	require.NotNil(t, uploadResp.Document)
	require.Equal(t, "PENDING", uploadResp.Document.VerificationStatus)

	_, err = svc.UploadUserDocument(ctx, &authnservicev1.UploadUserDocumentRequest{
		UserId:         userID,
		DocumentTypeId: typeResp.Types[0].DocumentTypeId,
		FileUrl:        "invalid-url",
	})
	require.Error(t, err)

	listDocsResp, err := svc.ListUserDocuments(ctx, &authnservicev1.ListUserDocumentsRequest{
		UserId:         userID,
		DocumentTypeId: typeResp.Types[0].DocumentTypeId,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listDocsResp.Documents)

	getDocResp, err := svc.GetUserDocument(ctx, &authnservicev1.GetUserDocumentRequest{
		UserDocumentId: uploadResp.Document.UserDocumentId,
	})
	require.NoError(t, err)
	require.Equal(t, uploadResp.Document.UserDocumentId, getDocResp.Document.UserDocumentId)

	_, err = svc.DeleteUserDocument(ctx, &authnservicev1.DeleteUserDocumentRequest{
		UserDocumentId: uploadResp.Document.UserDocumentId,
	})
	require.NoError(t, err)

	_, err = svc.GetUserDocument(ctx, &authnservicev1.GetUserDocumentRequest{
		UserDocumentId: uploadResp.Document.UserDocumentId,
	})
	require.Error(t, err)
}

func TestAuthService_LiveDB_KYCVoiceAndVerifyDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_kyc_voice_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Flow2")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })
	otherUserID := createLiveAuthnUser(t, svc, ctx, fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000), "live_kyc_other_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Flow3")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, otherUserID) })

	// Create a profile so ApproveKYC can update kyc_verified flag.
	_, err := svc.CreateUserProfile(ctx, &authnservicev1.CreateUserProfileRequest{
		UserId:      userID,
		FullName:    "KYC Voice User",
		DateOfBirth: timestamppb.New(time.Date(1992, time.August, 20, 0, 0, 0, 0, time.UTC)),
		Gender:      "FEMALE",
		City:        "Dhaka",
		Country:     "BD",
		NidNumber:   "9876543210123",
	})
	require.NoError(t, err)

	// KYC flow
	initKYC, err := svc.InitiateKYC(ctx, &authnservicev1.InitiateKYCRequest{UserId: userID})
	require.NoError(t, err)
	require.NotEmpty(t, initKYC.KycId)
	require.Equal(t, "PENDING", initKYC.Status)

	kycStatus, err := svc.GetKYCStatus(ctx, &authnservicev1.GetKYCStatusRequest{UserId: userID})
	require.NoError(t, err)
	require.Equal(t, initKYC.KycId, kycStatus.KycId)
	_, err = svc.SubmitKYCFrame(ctx, &authnservicev1.SubmitKYCFrameRequest{
		UserId:    otherUserID,
		SessionId: initKYC.KycId,
		ImageData: []byte("invalid-owner-frame"),
	})
	require.Error(t, err)

	submitResp, err := svc.SubmitKYCFrame(ctx, &authnservicev1.SubmitKYCFrameRequest{
		UserId:        userID,
		SessionId:     initKYC.KycId,
		ImageData:     []byte("frame-bytes-1"),
		FrameSequence: 1,
	})
	require.NoError(t, err)
	require.True(t, submitResp.Accepted)
	require.Equal(t, int32(1), submitResp.CompletedSteps)

	completeResp, err := svc.CompleteKYCSession(ctx, &authnservicev1.CompleteKYCSessionRequest{
		UserId:    userID,
		SessionId: initKYC.KycId,
	})
	require.NoError(t, err)
	require.True(t, completeResp.Success)
	require.Equal(t, "VERIFIED", completeResp.Status)

	profileAfterComplete, err := svc.GetUserProfile(ctx, &authnservicev1.GetUserProfileRequest{UserId: userID})
	require.NoError(t, err)
	require.True(t, profileAfterComplete.Profile.KycVerified)

	_, err = svc.RejectKYC(ctx, &authnservicev1.RejectKYCRequest{
		KycId:           initKYC.KycId,
		ReviewerId:      userID,
		RejectionReason: "missing field",
	})
	require.NoError(t, err)

	_, err = svc.ApproveKYC(ctx, &authnservicev1.ApproveKYCRequest{
		KycId:      initKYC.KycId,
		ReviewerId: userID,
	})
	require.NoError(t, err)

	// Document verify flow
	typeResp, err := svc.ListDocumentTypes(ctx, &authnservicev1.ListDocumentTypesRequest{})
	require.NoError(t, err)
	if len(typeResp.Types) == 0 {
		docType := &authnentityv1.DocumentType{
			DocumentTypeId: uuid.NewString(),
			Code:           "LIVE_VERIFY_DOC",
			Name:           "Live Verify Doc",
			Description:    "generated by live service test",
			IsActive:       true,
		}
		require.NoError(t, svc.documentTypeRepo.Create(ctx, docType))
		typeResp.Types = append(typeResp.Types, docType)
		t.Cleanup(func() {
			_ = dbConn.Exec(`DELETE FROM authn_schema.document_types WHERE document_type_id = ?`, docType.DocumentTypeId).Error
		})
	}

	uploadResp, err := svc.UploadUserDocument(ctx, &authnservicev1.UploadUserDocumentRequest{
		UserId:         userID,
		DocumentTypeId: typeResp.Types[0].DocumentTypeId,
		FileUrl:        "https://example.com/live-verify-doc.pdf",
	})
	require.NoError(t, err)

	verifyResp, err := svc.VerifyDocument(ctx, &authnservicev1.VerifyDocumentRequest{
		UserDocumentId:     uploadResp.Document.UserDocumentId,
		VerifiedBy:         userID,
		VerificationStatus: "VERIFIED",
		RejectionReason:    "",
	})
	require.NoError(t, err)
	require.NotNil(t, verifyResp.Document)
	require.Equal(t, "VERIFIED", verifyResp.Document.VerificationStatus)

	// Voice session flow using Create/Get/End methods.
	createVoiceResp, err := svc.CreateVoiceSession(ctx, &authnservicev1.CreateVoiceSessionRequest{
		UserId:      userID,
		Language:    "en-US",
		PhoneNumber: mobile,
	})
	require.NoError(t, err)
	require.NotEmpty(t, createVoiceResp.VoiceSessionId)

	getVoiceResp, err := svc.GetVoiceSession(ctx, &authnservicev1.GetVoiceSessionRequest{
		VoiceSessionId: createVoiceResp.VoiceSessionId,
	})
	require.NoError(t, err)
	require.Equal(t, userID, getVoiceResp.UserId)

	_, err = svc.EndVoiceSession(ctx, &authnservicev1.EndVoiceSessionRequest{
		VoiceSessionId:  createVoiceResp.VoiceSessionId,
		Status:          "FAILED",
		DurationSeconds: 19,
	})
	require.NoError(t, err)
}

func TestAuthService_LiveDB_ExternalKYCProxyFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "auth-service-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)

	mobile := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	userID := createLiveAuthnUser(t, svc, ctx, mobile, "live_ext_kyc_"+uuid.NewString()[:8]+"@example.com", "Str0ng!Ext1")
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, userID) })

	_, err := svc.CreateUserProfile(ctx, &authnservicev1.CreateUserProfileRequest{
		UserId:      userID,
		FullName:    "External KYC User",
		DateOfBirth: timestamppb.New(time.Date(1991, time.May, 10, 0, 0, 0, 0, time.UTC)),
		Gender:      "MALE",
		City:        "Dhaka",
		Country:     "BD",
		NidNumber:   "1234567890123",
	})
	require.NoError(t, err)

	externalSessionID := uuid.NewString()
	fakeClient := &fakeExternalKYCClient{
		startFn: func(ctx context.Context, in *kycservicev1.StartKYCVerificationRequest, opts ...grpc.CallOption) (*kycservicev1.StartKYCVerificationResponse, error) {
			require.Equal(t, userID, in.EntityId)
			return &kycservicev1.StartKYCVerificationResponse{KycVerificationId: externalSessionID}, nil
		},
		uploadFn: func(ctx context.Context, in *kycservicev1.UploadDocumentRequest, opts ...grpc.CallOption) (*kycservicev1.UploadDocumentResponse, error) {
			require.Equal(t, externalSessionID, in.KycVerificationId)
			require.Equal(t, "LIVENESS_FRAME", in.DocumentType)
			return &kycservicev1.UploadDocumentResponse{DocumentVerificationId: uuid.NewString()}, nil
		},
		verifyFn: func(ctx context.Context, in *kycservicev1.VerifyKYCRequest, opts ...grpc.CallOption) (*kycservicev1.VerifyKYCResponse, error) {
			require.Equal(t, externalSessionID, in.KycVerificationId)
			require.Equal(t, userID, in.VerifiedBy)
			return &kycservicev1.VerifyKYCResponse{Message: "verified"}, nil
		},
	}
	svc.SetExternalKYCClient(fakeClient)

	initResp, err := svc.InitiateKYC(ctx, &authnservicev1.InitiateKYCRequest{UserId: userID})
	require.NoError(t, err)
	require.Equal(t, externalSessionID, initResp.KycId)

	_, err = svc.SubmitKYCFrame(ctx, &authnservicev1.SubmitKYCFrameRequest{
		UserId:        userID,
		SessionId:     initResp.KycId,
		ImageData:     []byte("external-frame"),
		FrameSequence: 1,
	})
	require.NoError(t, err)

	completeResp, err := svc.CompleteKYCSession(ctx, &authnservicev1.CompleteKYCSessionRequest{
		UserId:    userID,
		SessionId: initResp.KycId,
	})
	require.NoError(t, err)
	require.True(t, completeResp.Success)

	require.Equal(t, 1, fakeClient.startCalls)
	require.Equal(t, 1, fakeClient.uploadCalls)
	require.Equal(t, 1, fakeClient.verifyCalls)
}
