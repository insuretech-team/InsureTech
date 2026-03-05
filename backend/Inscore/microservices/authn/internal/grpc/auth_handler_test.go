package grpc_test

import (
	"context"
	"errors"
	"testing"

	grpchandler "github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/grpc"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─────────────────────────────────────────────
// Minimal mock that satisfies AuthServiceIface
// ─────────────────────────────────────────────

type mockAuthService struct {
	registerErr                    error
	loginErr                       error
	sendOTPErr                     error
	verifyOTPErr                   error
	resendOTPErr                   error
	refreshTokenErr                error
	logoutErr                      error
	changePasswordErr              error
	resetPasswordErr               error
	validateTokenErr               error
	getSessionErr                  error
	listSessionsErr                error
	revokeSessionErr               error
	validateCSRFErr                error
	getCurrentSessionErr           error
	revokeAllSessionsErr           error
	registerEmailUserErr           error
	sendEmailOTPErr                error
	verifyEmailErr                 error
	emailLoginErr                  error
	requestPasswordResetByEmailErr error
	resetPasswordByEmailErr        error
	biometricAuthenticateErr       error
	updateDLRStatusErr             error
	createAPIKeyErr                error
	listAPIKeysErr                 error
	revokeAPIKeyErr                error
	rotateAPIKeyErr                error
	// profile / document
	createUserProfileErr  error
	getUserProfileErr     error
	updateUserProfileErr  error
	uploadUserDocumentErr error
	getUserDocumentErr    error
	updateUserDocumentErr error
	listUserDocumentsErr  error
	deleteUserDocumentErr error
	listDocumentTypesErr  error
	submitKYCFrameErr     error
	completeKYCSessionErr error
}

func (m *mockAuthService) Register(_ context.Context, _ *authnservicev1.RegisterRequest) (*authnservicev1.RegisterResponse, error) {
	return &authnservicev1.RegisterResponse{}, m.registerErr
}
func (m *mockAuthService) Login(_ context.Context, _ *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error) {
	return &authnservicev1.LoginResponse{}, m.loginErr
}
func (m *mockAuthService) SendOTP(_ context.Context, _ *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error) {
	return &authnservicev1.SendOTPResponse{}, m.sendOTPErr
}
func (m *mockAuthService) VerifyOTP(_ context.Context, _ *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error) {
	return &authnservicev1.VerifyOTPResponse{}, m.verifyOTPErr
}
func (m *mockAuthService) ResendOTP(_ context.Context, _ *authnservicev1.ResendOTPRequest) (*authnservicev1.ResendOTPResponse, error) {
	return &authnservicev1.ResendOTPResponse{}, m.resendOTPErr
}
func (m *mockAuthService) RefreshToken(_ context.Context, _ *authnservicev1.RefreshTokenRequest) (*authnservicev1.RefreshTokenResponse, error) {
	return &authnservicev1.RefreshTokenResponse{}, m.refreshTokenErr
}
func (m *mockAuthService) Logout(_ context.Context, _ *authnservicev1.LogoutRequest) (*authnservicev1.LogoutResponse, error) {
	return &authnservicev1.LogoutResponse{}, m.logoutErr
}
func (m *mockAuthService) ChangePassword(_ context.Context, _ *authnservicev1.ChangePasswordRequest) (*authnservicev1.ChangePasswordResponse, error) {
	return &authnservicev1.ChangePasswordResponse{}, m.changePasswordErr
}
func (m *mockAuthService) ResetPassword(_ context.Context, _ *authnservicev1.ResetPasswordRequest) (*authnservicev1.ResetPasswordResponse, error) {
	return &authnservicev1.ResetPasswordResponse{}, m.resetPasswordErr
}
func (m *mockAuthService) ValidateToken(_ context.Context, _ *authnservicev1.ValidateTokenRequest) (*authnservicev1.ValidateTokenResponse, error) {
	return &authnservicev1.ValidateTokenResponse{}, m.validateTokenErr
}
func (m *mockAuthService) GetSession(_ context.Context, _ *authnservicev1.GetSessionRequest) (*authnservicev1.GetSessionResponse, error) {
	return &authnservicev1.GetSessionResponse{}, m.getSessionErr
}
func (m *mockAuthService) ListSessions(_ context.Context, _ *authnservicev1.ListSessionsRequest) (*authnservicev1.ListSessionsResponse, error) {
	return &authnservicev1.ListSessionsResponse{}, m.listSessionsErr
}
func (m *mockAuthService) RevokeSession(_ context.Context, _ *authnservicev1.RevokeSessionRequest) (*authnservicev1.RevokeSessionResponse, error) {
	return &authnservicev1.RevokeSessionResponse{}, m.revokeSessionErr
}
func (m *mockAuthService) ValidateCSRF(_ context.Context, _ *authnservicev1.ValidateCSRFRequest) (*authnservicev1.ValidateCSRFResponse, error) {
	return &authnservicev1.ValidateCSRFResponse{}, m.validateCSRFErr
}
func (m *mockAuthService) GetCurrentSession(_ context.Context, _ *authnservicev1.GetCurrentSessionRequest) (*authnservicev1.GetCurrentSessionResponse, error) {
	return &authnservicev1.GetCurrentSessionResponse{}, m.getCurrentSessionErr
}
func (m *mockAuthService) RevokeAllSessions(_ context.Context, _ *authnservicev1.RevokeAllSessionsRequest) (*authnservicev1.RevokeAllSessionsResponse, error) {
	return &authnservicev1.RevokeAllSessionsResponse{}, m.revokeAllSessionsErr
}
func (m *mockAuthService) RegisterEmailUser(_ context.Context, _ *authnservicev1.RegisterEmailUserRequest) (*authnservicev1.RegisterEmailUserResponse, error) {
	return &authnservicev1.RegisterEmailUserResponse{}, m.registerEmailUserErr
}
func (m *mockAuthService) SendEmailOTP(_ context.Context, _ *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error) {
	return &authnservicev1.SendEmailOTPResponse{}, m.sendEmailOTPErr
}
func (m *mockAuthService) VerifyEmail(_ context.Context, _ *authnservicev1.VerifyEmailRequest) (*authnservicev1.VerifyEmailResponse, error) {
	return &authnservicev1.VerifyEmailResponse{}, m.verifyEmailErr
}
func (m *mockAuthService) EmailLogin(_ context.Context, _ *authnservicev1.EmailLoginRequest) (*authnservicev1.EmailLoginResponse, error) {
	return &authnservicev1.EmailLoginResponse{}, m.emailLoginErr
}
func (m *mockAuthService) RequestPasswordResetByEmail(_ context.Context, _ *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
	return &authnservicev1.RequestPasswordResetByEmailResponse{}, m.requestPasswordResetByEmailErr
}
func (m *mockAuthService) ResetPasswordByEmail(_ context.Context, _ *authnservicev1.ResetPasswordByEmailRequest) (*authnservicev1.ResetPasswordByEmailResponse, error) {
	return &authnservicev1.ResetPasswordByEmailResponse{}, m.resetPasswordByEmailErr
}
func (m *mockAuthService) BiometricAuthenticate(_ context.Context, _ *authnservicev1.BiometricAuthenticateRequest) (*authnservicev1.BiometricAuthenticateResponse, error) {
	return &authnservicev1.BiometricAuthenticateResponse{}, m.biometricAuthenticateErr
}
func (m *mockAuthService) UpdateDLRStatus(_ context.Context, _ *authnservicev1.UpdateDLRStatusRequest) (*authnservicev1.UpdateDLRStatusResponse, error) {
	return &authnservicev1.UpdateDLRStatusResponse{}, m.updateDLRStatusErr
}
func (m *mockAuthService) CreateAPIKey(_ context.Context, _ *authnservicev1.CreateAPIKeyRequest) (*authnservicev1.CreateAPIKeyResponse, error) {
	return &authnservicev1.CreateAPIKeyResponse{}, m.createAPIKeyErr
}
func (m *mockAuthService) ListAPIKeys(_ context.Context, _ *authnservicev1.ListAPIKeysRequest) (*authnservicev1.ListAPIKeysResponse, error) {
	return &authnservicev1.ListAPIKeysResponse{}, m.listAPIKeysErr
}
func (m *mockAuthService) RevokeAPIKey(_ context.Context, _ *authnservicev1.RevokeAPIKeyRequest) (*authnservicev1.RevokeAPIKeyResponse, error) {
	return &authnservicev1.RevokeAPIKeyResponse{}, m.revokeAPIKeyErr
}
func (m *mockAuthService) RotateAPIKey(_ context.Context, _ *authnservicev1.RotateAPIKeyRequest) (*authnservicev1.RotateAPIKeyResponse, error) {
	return &authnservicev1.RotateAPIKeyResponse{}, m.rotateAPIKeyErr
}
func (m *mockAuthService) CreateUserProfile(_ context.Context, _ *authnservicev1.CreateUserProfileRequest) (*authnservicev1.CreateUserProfileResponse, error) {
	return &authnservicev1.CreateUserProfileResponse{}, m.createUserProfileErr
}
func (m *mockAuthService) GetUserProfile(_ context.Context, _ *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	return &authnservicev1.GetUserProfileResponse{}, m.getUserProfileErr
}
func (m *mockAuthService) UpdateUserProfile(_ context.Context, _ *authnservicev1.UpdateUserProfileRequest) (*authnservicev1.UpdateUserProfileResponse, error) {
	return &authnservicev1.UpdateUserProfileResponse{}, m.updateUserProfileErr
}
func (m *mockAuthService) UploadUserDocument(_ context.Context, _ *authnservicev1.UploadUserDocumentRequest) (*authnservicev1.UploadUserDocumentResponse, error) {
	return &authnservicev1.UploadUserDocumentResponse{}, m.uploadUserDocumentErr
}
func (m *mockAuthService) GetUserDocument(_ context.Context, _ *authnservicev1.GetUserDocumentRequest) (*authnservicev1.GetUserDocumentResponse, error) {
	return &authnservicev1.GetUserDocumentResponse{}, m.getUserDocumentErr
}
func (m *mockAuthService) UpdateUserDocument(_ context.Context, _ *authnservicev1.UpdateUserDocumentRequest) (*authnservicev1.UpdateUserDocumentResponse, error) {
	return &authnservicev1.UpdateUserDocumentResponse{}, m.updateUserDocumentErr
}
func (m *mockAuthService) ListUserDocuments(_ context.Context, _ *authnservicev1.ListUserDocumentsRequest) (*authnservicev1.ListUserDocumentsResponse, error) {
	return &authnservicev1.ListUserDocumentsResponse{}, m.listUserDocumentsErr
}
func (m *mockAuthService) DeleteUserDocument(_ context.Context, _ *authnservicev1.DeleteUserDocumentRequest) (*authnservicev1.DeleteUserDocumentResponse, error) {
	return &authnservicev1.DeleteUserDocumentResponse{}, m.deleteUserDocumentErr
}
func (m *mockAuthService) ListDocumentTypes(_ context.Context, _ *authnservicev1.ListDocumentTypesRequest) (*authnservicev1.ListDocumentTypesResponse, error) {
	return &authnservicev1.ListDocumentTypesResponse{}, m.listDocumentTypesErr
}

// ── Profile / Document / KYC / Voice / TOTP stubs for the mock ────────────────
func (m *mockAuthService) InitiateKYC(_ context.Context, _ *authnservicev1.InitiateKYCRequest) (*authnservicev1.InitiateKYCResponse, error) {
	return &authnservicev1.InitiateKYCResponse{}, nil
}
func (m *mockAuthService) GetKYCStatus(_ context.Context, _ *authnservicev1.GetKYCStatusRequest) (*authnservicev1.GetKYCStatusResponse, error) {
	return &authnservicev1.GetKYCStatusResponse{}, nil
}
func (m *mockAuthService) SubmitKYCFrame(_ context.Context, _ *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	return &authnservicev1.SubmitKYCFrameResponse{}, m.submitKYCFrameErr
}
func (m *mockAuthService) CompleteKYCSession(_ context.Context, _ *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
	return &authnservicev1.CompleteKYCSessionResponse{}, m.completeKYCSessionErr
}
func (m *mockAuthService) ApproveKYC(_ context.Context, _ *authnservicev1.ApproveKYCRequest) (*authnservicev1.ApproveKYCResponse, error) {
	return &authnservicev1.ApproveKYCResponse{}, nil
}
func (m *mockAuthService) RejectKYC(_ context.Context, _ *authnservicev1.RejectKYCRequest) (*authnservicev1.RejectKYCResponse, error) {
	return &authnservicev1.RejectKYCResponse{}, nil
}
func (m *mockAuthService) VerifyDocument(_ context.Context, _ *authnservicev1.VerifyDocumentRequest) (*authnservicev1.VerifyDocumentResponse, error) {
	return &authnservicev1.VerifyDocumentResponse{}, nil
}
func (m *mockAuthService) CreateVoiceSession(_ context.Context, _ *authnservicev1.CreateVoiceSessionRequest) (*authnservicev1.CreateVoiceSessionResponse, error) {
	return &authnservicev1.CreateVoiceSessionResponse{}, nil
}
func (m *mockAuthService) GetVoiceSession(_ context.Context, _ *authnservicev1.GetVoiceSessionRequest) (*authnservicev1.GetVoiceSessionResponse, error) {
	return &authnservicev1.GetVoiceSessionResponse{}, nil
}
func (m *mockAuthService) EndVoiceSession(_ context.Context, _ *authnservicev1.EndVoiceSessionRequest) (*authnservicev1.EndVoiceSessionResponse, error) {
	return &authnservicev1.EndVoiceSessionResponse{}, nil
}
func (m *mockAuthService) GetProfilePhotoUploadURL(_ context.Context, _ *authnservicev1.GetProfilePhotoUploadURLRequest) (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
	return &authnservicev1.GetProfilePhotoUploadURLResponse{}, nil
}
func (m *mockAuthService) UpdateNotificationPreferences(_ context.Context, _ *authnservicev1.UpdateNotificationPreferencesRequest) (*authnservicev1.UpdateNotificationPreferencesResponse, error) {
	return &authnservicev1.UpdateNotificationPreferencesResponse{}, nil
}
func (m *mockAuthService) EnableTOTP(_ context.Context, _ *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	return &authnservicev1.EnableTOTPResponse{}, nil
}
func (m *mockAuthService) VerifyTOTP(_ context.Context, _ *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	return &authnservicev1.VerifyTOTPResponse{}, nil
}
func (m *mockAuthService) DisableTOTP(_ context.Context, _ *authnservicev1.DisableTOTPRequest) (*authnservicev1.DisableTOTPResponse, error) {
	return &authnservicev1.DisableTOTPResponse{}, nil
}
func (m *mockAuthService) GetJWKS(_ context.Context, _ *authnservicev1.GetJWKSRequest) (*authnservicev1.GetJWKSResponse, error) {
	return &authnservicev1.GetJWKSResponse{}, nil
}
func (m *mockAuthService) InitiateVoiceSession(_ context.Context, _ *authnservicev1.InitiateVoiceSessionRequest) (*authnservicev1.InitiateVoiceSessionResponse, error) {
	return &authnservicev1.InitiateVoiceSessionResponse{}, nil
}
func (m *mockAuthService) SubmitVoiceSample(_ context.Context, _ *authnservicev1.SubmitVoiceSampleRequest) (*authnservicev1.SubmitVoiceSampleResponse, error) {
	return &authnservicev1.SubmitVoiceSampleResponse{}, nil
}
func (m *mockAuthService) VerifyVoiceSession(_ context.Context, _ *authnservicev1.VerifyVoiceSessionRequest) (*authnservicev1.VerifyVoiceSessionResponse, error) {
	return &authnservicev1.VerifyVoiceSessionResponse{}, nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func newHandler(svc *mockAuthService) *grpchandler.AuthServiceHandler {
	return grpchandler.NewAuthServiceHandler(svc)
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", want)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error is not a gRPC status error: %v", err)
	}
	if st.Code() != want {
		t.Fatalf("expected gRPC code %v, got %v (msg: %s)", want, st.Code(), st.Message())
	}
}

// ─────────────────────────────────────────────
// Tests: input validation
// ─────────────────────────────────────────────

func TestRegister_MissingMobileNumber(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.Register(context.Background(), &authnservicev1.RegisterRequest{
		Password: "Secret123!",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestRegister_MissingPassword(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.Register(context.Background(), &authnservicev1.RegisterRequest{
		MobileNumber: "8801712345678",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestLogin_MissingMobileNumber(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.Login(context.Background(), &authnservicev1.LoginRequest{
		Password: "Secret123!",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestLogin_MissingPassword(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.Login(context.Background(), &authnservicev1.LoginRequest{
		MobileNumber: "8801712345678",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestValidateToken_MissingToken(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.ValidateToken(context.Background(), &authnservicev1.ValidateTokenRequest{})
	assertCode(t, err, codes.InvalidArgument)
}

func TestVerifyOTP_MissingOtpId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.VerifyOTP(context.Background(), &authnservicev1.VerifyOTPRequest{Code: "123456"})
	assertCode(t, err, codes.InvalidArgument)
}

func TestVerifyOTP_MissingCode(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.VerifyOTP(context.Background(), &authnservicev1.VerifyOTPRequest{OtpId: "uuid"})
	assertCode(t, err, codes.InvalidArgument)
}

func TestBiometricAuthenticate_MissingToken(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.BiometricAuthenticate(context.Background(), &authnservicev1.BiometricAuthenticateRequest{
		DeviceId: "dev-1",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestBiometricAuthenticate_MissingDeviceId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.BiometricAuthenticate(context.Background(), &authnservicev1.BiometricAuthenticateRequest{
		BiometricToken: "tok",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestCreateAPIKey_MissingOwnerId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.CreateAPIKey(context.Background(), &authnservicev1.CreateAPIKeyRequest{Name: "my-key"})
	assertCode(t, err, codes.InvalidArgument)
}

func TestCreateAPIKey_MissingName(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.CreateAPIKey(context.Background(), &authnservicev1.CreateAPIKeyRequest{OwnerId: "user-1"})
	assertCode(t, err, codes.InvalidArgument)
}

func TestListAPIKeys_MissingOwnerId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.ListAPIKeys(context.Background(), &authnservicev1.ListAPIKeysRequest{})
	assertCode(t, err, codes.InvalidArgument)
}

func TestRevokeAPIKey_MissingKeyId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.RevokeAPIKey(context.Background(), &authnservicev1.RevokeAPIKeyRequest{})
	assertCode(t, err, codes.InvalidArgument)
}

// ─────────────────────────────────────────────
// Tests: error code mapping via MapError
// ─────────────────────────────────────────────

func TestMapError_NotFound(t *testing.T) {
	err := grpchandler.MapError(errors.New("user not found"))
	assertCode(t, err, codes.NotFound)
}

func TestMapError_AlreadyExists(t *testing.T) {
	err := grpchandler.MapError(errors.New("user already exists"))
	assertCode(t, err, codes.AlreadyExists)
}

func TestMapError_Unauthenticated(t *testing.T) {
	err := grpchandler.MapError(errors.New("invalid credentials"))
	assertCode(t, err, codes.Unauthenticated)
}

func TestMapError_PermissionDenied(t *testing.T) {
	err := grpchandler.MapError(errors.New("access denied"))
	assertCode(t, err, codes.PermissionDenied)
}

func TestMapError_ResourceExhausted(t *testing.T) {
	err := grpchandler.MapError(errors.New("too many requests"))
	assertCode(t, err, codes.ResourceExhausted)
}

func TestMapError_InvalidArgument(t *testing.T) {
	err := grpchandler.MapError(errors.New("invalid otp code"))
	assertCode(t, err, codes.Unauthenticated)
}

func TestMapError_Expired(t *testing.T) {
	// "otp expired" contains "expired" → Unauthenticated (session/token expired = auth failure)
	err := grpchandler.MapError(errors.New("otp expired"))
	assertCode(t, err, codes.Unauthenticated)
}

func TestMapError_Internal(t *testing.T) {
	err := grpchandler.MapError(errors.New("database connection failed"))
	assertCode(t, err, codes.Internal)
}

func TestMapError_NilPassthrough(t *testing.T) {
	// MapError(nil) should return nil — handler doesn't call mapError on nil
	// so this tests defensive nil handling if ever called directly.
	if grpchandler.MapError(nil) != nil {
		t.Fatal("expected nil for nil error input")
	}
}

// ─────────────────────────────────────────────
// Tests: service error propagation → correct gRPC code
// ─────────────────────────────────────────────

func TestRegister_ServiceNotFoundError(t *testing.T) {
	h := newHandler(&mockAuthService{registerErr: errors.New("user already exists")})
	_, err := h.Register(context.Background(), &authnservicev1.RegisterRequest{
		MobileNumber: "8801712345678",
		Password:     "Secret123!",
	})
	assertCode(t, err, codes.AlreadyExists)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := newHandler(&mockAuthService{loginErr: errors.New("invalid credentials")})
	_, err := h.Login(context.Background(), &authnservicev1.LoginRequest{
		MobileNumber: "8801712345678",
		Password:     "wrongpass",
	})
	assertCode(t, err, codes.Unauthenticated)
}

func TestRefreshToken_Expired(t *testing.T) {
	h := newHandler(&mockAuthService{refreshTokenErr: errors.New("token expired")})
	_, err := h.RefreshToken(context.Background(), &authnservicev1.RefreshTokenRequest{
		RefreshToken: "tok",
	})
	assertCode(t, err, codes.Unauthenticated)
}

func TestBiometricAuthenticate_Success(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.BiometricAuthenticate(context.Background(), &authnservicev1.BiometricAuthenticateRequest{
		BiometricToken: "valid-token",
		DeviceId:       "device-1",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestUpdateDLRStatus_MissingProviderMessageId(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.UpdateDLRStatus(context.Background(), &authnservicev1.UpdateDLRStatusRequest{
		Status: "DELIVERED",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestUpdateDLRStatus_MissingStatus(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.UpdateDLRStatus(context.Background(), &authnservicev1.UpdateDLRStatusRequest{
		ProviderMessageId: "msg-123",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestCreateAPIKey_Success(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.CreateAPIKey(context.Background(), &authnservicev1.CreateAPIKeyRequest{
		OwnerId: "user-1",
		Name:    "my-key",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestRevokeAPIKey_NotFound(t *testing.T) {
	h := newHandler(&mockAuthService{revokeAPIKeyErr: errors.New("key not found")})
	_, err := h.RevokeAPIKey(context.Background(), &authnservicev1.RevokeAPIKeyRequest{
		KeyId: "key-1",
	})
	assertCode(t, err, codes.NotFound)
}

func TestSubmitKYCFrame_MissingUserID(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.SubmitKYCFrame(context.Background(), &authnservicev1.SubmitKYCFrameRequest{
		SessionId: "s1",
		ImageData: []byte("img"),
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestSubmitKYCFrame_MissingSessionID(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.SubmitKYCFrame(context.Background(), &authnservicev1.SubmitKYCFrameRequest{
		UserId:    "u1",
		ImageData: []byte("img"),
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestSubmitKYCFrame_MissingImageData(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.SubmitKYCFrame(context.Background(), &authnservicev1.SubmitKYCFrameRequest{
		UserId:    "u1",
		SessionId: "s1",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestCompleteKYCSession_MissingUserID(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.CompleteKYCSession(context.Background(), &authnservicev1.CompleteKYCSessionRequest{
		SessionId: "s1",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestCompleteKYCSession_MissingSessionID(t *testing.T) {
	h := newHandler(&mockAuthService{})
	_, err := h.CompleteKYCSession(context.Background(), &authnservicev1.CompleteKYCSessionRequest{
		UserId: "u1",
	})
	assertCode(t, err, codes.InvalidArgument)
}

func TestKYCFrameAndComplete_ServiceErrorsMapped(t *testing.T) {
	h := newHandler(&mockAuthService{
		submitKYCFrameErr:     errors.New("forbidden"),
		completeKYCSessionErr: errors.New("not found"),
	})
	_, err := h.SubmitKYCFrame(context.Background(), &authnservicev1.SubmitKYCFrameRequest{
		UserId:    "u1",
		SessionId: "s1",
		ImageData: []byte("img"),
	})
	assertCode(t, err, codes.PermissionDenied)

	_, err = h.CompleteKYCSession(context.Background(), &authnservicev1.CompleteKYCSessionRequest{
		UserId:    "u1",
		SessionId: "missing",
	})
	assertCode(t, err, codes.NotFound)
}

func TestAuthHandler_SuccessCoverage_AllMajorRPCs(t *testing.T) {
	h := newHandler(&mockAuthService{})
	ctx := context.Background()

	_, err := h.SendOTP(ctx, &authnservicev1.SendOTPRequest{Recipient: "+8801712345678", Type: "login"})
	if err != nil {
		t.Fatalf("SendOTP: %v", err)
	}
	_, err = h.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{OtpId: "otp-1", Code: "123456"})
	if err != nil {
		t.Fatalf("VerifyOTP: %v", err)
	}
	_, err = h.ResendOTP(ctx, &authnservicev1.ResendOTPRequest{OriginalOtpId: "otp-1"})
	if err != nil {
		t.Fatalf("ResendOTP: %v", err)
	}
	_, err = h.RefreshToken(ctx, &authnservicev1.RefreshTokenRequest{RefreshToken: "rt"})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	_, err = h.Logout(ctx, &authnservicev1.LogoutRequest{SessionId: "s1"})
	if err != nil {
		t.Fatalf("Logout: %v", err)
	}
	_, err = h.GetSession(ctx, &authnservicev1.GetSessionRequest{SessionId: "s1"})
	if err != nil {
		t.Fatalf("GetSession: %v", err)
	}
	_, err = h.ListSessions(ctx, &authnservicev1.ListSessionsRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	_, err = h.RevokeSession(ctx, &authnservicev1.RevokeSessionRequest{SessionId: "s1"})
	if err != nil {
		t.Fatalf("RevokeSession: %v", err)
	}
	_, err = h.ValidateCSRF(ctx, &authnservicev1.ValidateCSRFRequest{SessionId: "s1", CsrfToken: "c1"})
	if err != nil {
		t.Fatalf("ValidateCSRF: %v", err)
	}
	_, err = h.GetCurrentSession(ctx, &authnservicev1.GetCurrentSessionRequest{})
	if err != nil {
		t.Fatalf("GetCurrentSession: %v", err)
	}
	_, err = h.RevokeAllSessions(ctx, &authnservicev1.RevokeAllSessionsRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("RevokeAllSessions: %v", err)
	}
	_, err = h.ChangePassword(ctx, &authnservicev1.ChangePasswordRequest{UserId: "u1", OldPassword: "Old123!@", NewPassword: "New123!@"})
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}
	_, err = h.ResetPassword(ctx, &authnservicev1.ResetPasswordRequest{MobileNumber: "+8801712345678", OtpCode: "123456", NewPassword: "New123!@"})
	if err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	_, err = h.RegisterEmailUser(ctx, &authnservicev1.RegisterEmailUserRequest{Email: "u@example.com", Password: "Pass123!@"})
	if err != nil {
		t.Fatalf("RegisterEmailUser: %v", err)
	}
	_, err = h.SendEmailOTP(ctx, &authnservicev1.SendEmailOTPRequest{Email: "u@example.com"})
	if err != nil {
		t.Fatalf("SendEmailOTP: %v", err)
	}
	_, err = h.VerifyEmail(ctx, &authnservicev1.VerifyEmailRequest{OtpId: "e1", Code: "123456"})
	if err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	_, err = h.EmailLogin(ctx, &authnservicev1.EmailLoginRequest{Email: "u@example.com", OtpId: "e1", Code: "123456"})
	if err != nil {
		t.Fatalf("EmailLogin: %v", err)
	}
	_, err = h.RequestPasswordResetByEmail(ctx, &authnservicev1.RequestPasswordResetByEmailRequest{Email: "u@example.com"})
	if err != nil {
		t.Fatalf("RequestPasswordResetByEmail: %v", err)
	}
	_, err = h.ResetPasswordByEmail(ctx, &authnservicev1.ResetPasswordByEmailRequest{Email: "u@example.com", OtpId: "e1", OtpCode: "123456", NewPassword: "New123!@"})
	if err != nil {
		t.Fatalf("ResetPasswordByEmail: %v", err)
	}
	_, err = h.UpdateDLRStatus(ctx, &authnservicev1.UpdateDLRStatusRequest{ProviderMessageId: "pm1", Status: "DELIVERED"})
	if err != nil {
		t.Fatalf("UpdateDLRStatus: %v", err)
	}
	_, err = h.ListAPIKeys(ctx, &authnservicev1.ListAPIKeysRequest{OwnerId: "u1"})
	if err != nil {
		t.Fatalf("ListAPIKeys: %v", err)
	}
	_, err = h.CreateUserProfile(ctx, &authnservicev1.CreateUserProfileRequest{
		UserId:       "u1",
		FullName:     "User One",
		AddressLine1: "Addr",
		City:         "Dhaka",
		District:     "Dhaka",
		Division:     "Dhaka",
		NidNumber:    "1234567890",
	})
	if err != nil {
		t.Fatalf("CreateUserProfile: %v", err)
	}
	_, err = h.GetUserProfile(ctx, &authnservicev1.GetUserProfileRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("GetUserProfile: %v", err)
	}
	_, err = h.UpdateUserProfile(ctx, &authnservicev1.UpdateUserProfileRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("UpdateUserProfile: %v", err)
	}
	_, err = h.UploadUserDocument(ctx, &authnservicev1.UploadUserDocumentRequest{UserId: "u1", DocumentTypeId: "dt1", FileUrl: "https://f"})
	if err != nil {
		t.Fatalf("UploadUserDocument: %v", err)
	}
	_, err = h.ListUserDocuments(ctx, &authnservicev1.ListUserDocumentsRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("ListUserDocuments: %v", err)
	}
	_, err = h.GetUserDocument(ctx, &authnservicev1.GetUserDocumentRequest{UserDocumentId: "d1"})
	if err != nil {
		t.Fatalf("GetUserDocument: %v", err)
	}
	_, err = h.UpdateUserDocument(ctx, &authnservicev1.UpdateUserDocumentRequest{UserDocumentId: "d1", FileUrl: "https://f2"})
	if err != nil {
		t.Fatalf("UpdateUserDocument: %v", err)
	}
	_, err = h.DeleteUserDocument(ctx, &authnservicev1.DeleteUserDocumentRequest{UserDocumentId: "d1"})
	if err != nil {
		t.Fatalf("DeleteUserDocument: %v", err)
	}
	_, err = h.ListDocumentTypes(ctx, &authnservicev1.ListDocumentTypesRequest{})
	if err != nil {
		t.Fatalf("ListDocumentTypes: %v", err)
	}
	_, err = h.InitiateKYC(ctx, &authnservicev1.InitiateKYCRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("InitiateKYC: %v", err)
	}
	_, err = h.GetKYCStatus(ctx, &authnservicev1.GetKYCStatusRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("GetKYCStatus: %v", err)
	}
	_, err = h.SubmitKYCFrame(ctx, &authnservicev1.SubmitKYCFrameRequest{UserId: "u1", SessionId: "k-session-1", ImageData: []byte("image-bytes")})
	if err != nil {
		t.Fatalf("SubmitKYCFrame: %v", err)
	}
	_, err = h.CompleteKYCSession(ctx, &authnservicev1.CompleteKYCSessionRequest{UserId: "u1", SessionId: "k-session-1"})
	if err != nil {
		t.Fatalf("CompleteKYCSession: %v", err)
	}
	_, err = h.ApproveKYC(ctx, &authnservicev1.ApproveKYCRequest{KycId: "k1"})
	if err != nil {
		t.Fatalf("ApproveKYC: %v", err)
	}
	_, err = h.RejectKYC(ctx, &authnservicev1.RejectKYCRequest{KycId: "k1", RejectionReason: "bad"})
	if err != nil {
		t.Fatalf("RejectKYC: %v", err)
	}
	_, err = h.VerifyDocument(ctx, &authnservicev1.VerifyDocumentRequest{UserDocumentId: "d1"})
	if err != nil {
		t.Fatalf("VerifyDocument: %v", err)
	}
	_, err = h.CreateVoiceSession(ctx, &authnservicev1.CreateVoiceSessionRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("CreateVoiceSession: %v", err)
	}
	_, err = h.GetVoiceSession(ctx, &authnservicev1.GetVoiceSessionRequest{VoiceSessionId: "vs1"})
	if err != nil {
		t.Fatalf("GetVoiceSession: %v", err)
	}
	_, err = h.EndVoiceSession(ctx, &authnservicev1.EndVoiceSessionRequest{VoiceSessionId: "vs1"})
	if err != nil {
		t.Fatalf("EndVoiceSession: %v", err)
	}
	_, err = h.InitiateVoiceSession(ctx, &authnservicev1.InitiateVoiceSessionRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("InitiateVoiceSession: %v", err)
	}
	_, err = h.SubmitVoiceSample(ctx, &authnservicev1.SubmitVoiceSampleRequest{SessionId: "vs1"})
	if err != nil {
		t.Fatalf("SubmitVoiceSample: %v", err)
	}
	_, err = h.VerifyVoiceSession(ctx, &authnservicev1.VerifyVoiceSessionRequest{SessionId: "vs1"})
	if err != nil {
		t.Fatalf("VerifyVoiceSession: %v", err)
	}
	_, err = h.GetProfilePhotoUploadURL(ctx, &authnservicev1.GetProfilePhotoUploadURLRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("GetProfilePhotoUploadURL: %v", err)
	}
	_, err = h.UpdateNotificationPreferences(ctx, &authnservicev1.UpdateNotificationPreferencesRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("UpdateNotificationPreferences: %v", err)
	}
	_, err = h.EnableTOTP(ctx, &authnservicev1.EnableTOTPRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("EnableTOTP: %v", err)
	}
	_, err = h.VerifyTOTP(ctx, &authnservicev1.VerifyTOTPRequest{UserId: "u1", TotpCode: "123456"})
	if err != nil {
		t.Fatalf("VerifyTOTP: %v", err)
	}
	_, err = h.DisableTOTP(ctx, &authnservicev1.DisableTOTPRequest{UserId: "u1"})
	if err != nil {
		t.Fatalf("DisableTOTP: %v", err)
	}
	_, err = h.GetJWKS(ctx, &authnservicev1.GetJWKSRequest{})
	if err != nil {
		t.Fatalf("GetJWKS: %v", err)
	}
}
