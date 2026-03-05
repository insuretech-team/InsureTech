package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockAuthServiceClient struct {
	// embed to satisfy forward-compat on server side (not required for client interface)
	authnservicev1.UnimplementedAuthServiceServer

	loginFn         func(ctx context.Context, in *authnservicev1.LoginRequest, opts ...grpc.CallOption) (*authnservicev1.LoginResponse, error)
	getSessionFn    func(ctx context.Context, in *authnservicev1.GetSessionRequest, opts ...grpc.CallOption) (*authnservicev1.GetSessionResponse, error)
	listSessionsFn  func(ctx context.Context, in *authnservicev1.ListSessionsRequest, opts ...grpc.CallOption) (*authnservicev1.ListSessionsResponse, error)
	revokeSessionFn func(ctx context.Context, in *authnservicev1.RevokeSessionRequest, opts ...grpc.CallOption) (*authnservicev1.RevokeSessionResponse, error)
	resendOTPFn     func(ctx context.Context, in *authnservicev1.ResendOTPRequest, opts ...grpc.CallOption) (*authnservicev1.ResendOTPResponse, error)
	verifyOTPFn     func(ctx context.Context, in *authnservicev1.VerifyOTPRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyOTPResponse, error)
	emailLoginFn    func(ctx context.Context, in *authnservicev1.EmailLoginRequest, opts ...grpc.CallOption) (*authnservicev1.EmailLoginResponse, error)
}

// Ensure we implement the generated gRPC client interface.
var _ authnservicev1.AuthServiceClient = (*mockAuthServiceClient)(nil)

func (m *mockAuthServiceClient) Register(ctx context.Context, in *authnservicev1.RegisterRequest, opts ...grpc.CallOption) (*authnservicev1.RegisterResponse, error) {
	return &authnservicev1.RegisterResponse{}, nil
}
func (m *mockAuthServiceClient) SendOTP(ctx context.Context, in *authnservicev1.SendOTPRequest, opts ...grpc.CallOption) (*authnservicev1.SendOTPResponse, error) {
	return &authnservicev1.SendOTPResponse{}, nil
}
func (m *mockAuthServiceClient) VerifyOTP(ctx context.Context, in *authnservicev1.VerifyOTPRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyOTPResponse, error) {
	if m.verifyOTPFn != nil {
		return m.verifyOTPFn(ctx, in, opts...)
	}
	return &authnservicev1.VerifyOTPResponse{}, nil
}
func (m *mockAuthServiceClient) ResendOTP(ctx context.Context, in *authnservicev1.ResendOTPRequest, opts ...grpc.CallOption) (*authnservicev1.ResendOTPResponse, error) {
	if m.resendOTPFn != nil {
		return m.resendOTPFn(ctx, in, opts...)
	}
	return &authnservicev1.ResendOTPResponse{}, nil
}
func (m *mockAuthServiceClient) Login(ctx context.Context, in *authnservicev1.LoginRequest, opts ...grpc.CallOption) (*authnservicev1.LoginResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, in, opts...)
	}
	return &authnservicev1.LoginResponse{}, nil
}
func (m *mockAuthServiceClient) RefreshToken(ctx context.Context, in *authnservicev1.RefreshTokenRequest, opts ...grpc.CallOption) (*authnservicev1.RefreshTokenResponse, error) {
	return &authnservicev1.RefreshTokenResponse{}, nil
}
func (m *mockAuthServiceClient) Logout(ctx context.Context, in *authnservicev1.LogoutRequest, opts ...grpc.CallOption) (*authnservicev1.LogoutResponse, error) {
	return &authnservicev1.LogoutResponse{}, nil
}
func (m *mockAuthServiceClient) ChangePassword(ctx context.Context, in *authnservicev1.ChangePasswordRequest, opts ...grpc.CallOption) (*authnservicev1.ChangePasswordResponse, error) {
	return &authnservicev1.ChangePasswordResponse{}, nil
}
func (m *mockAuthServiceClient) ResetPassword(ctx context.Context, in *authnservicev1.ResetPasswordRequest, opts ...grpc.CallOption) (*authnservicev1.ResetPasswordResponse, error) {
	return &authnservicev1.ResetPasswordResponse{}, nil
}
func (m *mockAuthServiceClient) ValidateToken(ctx context.Context, in *authnservicev1.ValidateTokenRequest, opts ...grpc.CallOption) (*authnservicev1.ValidateTokenResponse, error) {
	return &authnservicev1.ValidateTokenResponse{Valid: true}, nil
}
func (m *mockAuthServiceClient) GetSession(ctx context.Context, in *authnservicev1.GetSessionRequest, opts ...grpc.CallOption) (*authnservicev1.GetSessionResponse, error) {
	if m.getSessionFn != nil {
		return m.getSessionFn(ctx, in, opts...)
	}
	return &authnservicev1.GetSessionResponse{}, nil
}
func (m *mockAuthServiceClient) ListSessions(ctx context.Context, in *authnservicev1.ListSessionsRequest, opts ...grpc.CallOption) (*authnservicev1.ListSessionsResponse, error) {
	if m.listSessionsFn != nil {
		return m.listSessionsFn(ctx, in, opts...)
	}
	return &authnservicev1.ListSessionsResponse{}, nil
}
func (m *mockAuthServiceClient) RevokeSession(ctx context.Context, in *authnservicev1.RevokeSessionRequest, opts ...grpc.CallOption) (*authnservicev1.RevokeSessionResponse, error) {
	if m.revokeSessionFn != nil {
		return m.revokeSessionFn(ctx, in, opts...)
	}
	return &authnservicev1.RevokeSessionResponse{}, nil
}
func (m *mockAuthServiceClient) ValidateCSRF(ctx context.Context, in *authnservicev1.ValidateCSRFRequest, opts ...grpc.CallOption) (*authnservicev1.ValidateCSRFResponse, error) {
	return &authnservicev1.ValidateCSRFResponse{Valid: true}, nil
}
func (m *mockAuthServiceClient) GetCurrentSession(ctx context.Context, in *authnservicev1.GetCurrentSessionRequest, opts ...grpc.CallOption) (*authnservicev1.GetCurrentSessionResponse, error) {
	return &authnservicev1.GetCurrentSessionResponse{}, nil
}
func (m *mockAuthServiceClient) RevokeAllSessions(ctx context.Context, in *authnservicev1.RevokeAllSessionsRequest, opts ...grpc.CallOption) (*authnservicev1.RevokeAllSessionsResponse, error) {
	return &authnservicev1.RevokeAllSessionsResponse{}, nil
}
func (m *mockAuthServiceClient) RegisterEmailUser(ctx context.Context, in *authnservicev1.RegisterEmailUserRequest, opts ...grpc.CallOption) (*authnservicev1.RegisterEmailUserResponse, error) {
	return &authnservicev1.RegisterEmailUserResponse{}, nil
}
func (m *mockAuthServiceClient) SendEmailOTP(ctx context.Context, in *authnservicev1.SendEmailOTPRequest, opts ...grpc.CallOption) (*authnservicev1.SendEmailOTPResponse, error) {
	return &authnservicev1.SendEmailOTPResponse{}, nil
}
func (m *mockAuthServiceClient) VerifyEmail(ctx context.Context, in *authnservicev1.VerifyEmailRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyEmailResponse, error) {
	return &authnservicev1.VerifyEmailResponse{}, nil
}
func (m *mockAuthServiceClient) EmailLogin(ctx context.Context, in *authnservicev1.EmailLoginRequest, opts ...grpc.CallOption) (*authnservicev1.EmailLoginResponse, error) {
	if m.emailLoginFn != nil {
		return m.emailLoginFn(ctx, in, opts...)
	}
	return &authnservicev1.EmailLoginResponse{}, nil
}
func (m *mockAuthServiceClient) RequestPasswordResetByEmail(ctx context.Context, in *authnservicev1.RequestPasswordResetByEmailRequest, opts ...grpc.CallOption) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
	return &authnservicev1.RequestPasswordResetByEmailResponse{}, nil
}
func (m *mockAuthServiceClient) ResetPasswordByEmail(ctx context.Context, in *authnservicev1.ResetPasswordByEmailRequest, opts ...grpc.CallOption) (*authnservicev1.ResetPasswordByEmailResponse, error) {
	return &authnservicev1.ResetPasswordByEmailResponse{}, nil
}
func (m *mockAuthServiceClient) BiometricAuthenticate(ctx context.Context, in *authnservicev1.BiometricAuthenticateRequest, opts ...grpc.CallOption) (*authnservicev1.BiometricAuthenticateResponse, error) {
	return &authnservicev1.BiometricAuthenticateResponse{}, nil
}
func (m *mockAuthServiceClient) UpdateDLRStatus(ctx context.Context, in *authnservicev1.UpdateDLRStatusRequest, opts ...grpc.CallOption) (*authnservicev1.UpdateDLRStatusResponse, error) {
	return &authnservicev1.UpdateDLRStatusResponse{}, nil
}
func (m *mockAuthServiceClient) CreateAPIKey(ctx context.Context, in *authnservicev1.CreateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.CreateAPIKeyResponse, error) {
	return &authnservicev1.CreateAPIKeyResponse{}, nil
}
func (m *mockAuthServiceClient) ListAPIKeys(ctx context.Context, in *authnservicev1.ListAPIKeysRequest, opts ...grpc.CallOption) (*authnservicev1.ListAPIKeysResponse, error) {
	return &authnservicev1.ListAPIKeysResponse{}, nil
}
func (m *mockAuthServiceClient) RevokeAPIKey(ctx context.Context, in *authnservicev1.RevokeAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.RevokeAPIKeyResponse, error) {
	return &authnservicev1.RevokeAPIKeyResponse{}, nil
}
func (m *mockAuthServiceClient) RotateAPIKey(ctx context.Context, in *authnservicev1.RotateAPIKeyRequest, opts ...grpc.CallOption) (*authnservicev1.RotateAPIKeyResponse, error) {
	return &authnservicev1.RotateAPIKeyResponse{}, nil
}

// ── Profile / Document / KYC / Voice / TOTP stubs for the mock ────────────────
func (m *mockAuthServiceClient) CreateUserProfile(ctx context.Context, in *authnservicev1.CreateUserProfileRequest, opts ...grpc.CallOption) (*authnservicev1.CreateUserProfileResponse, error) {
	return &authnservicev1.CreateUserProfileResponse{}, nil
}
func (m *mockAuthServiceClient) GetUserProfile(ctx context.Context, in *authnservicev1.GetUserProfileRequest, opts ...grpc.CallOption) (*authnservicev1.GetUserProfileResponse, error) {
	return &authnservicev1.GetUserProfileResponse{}, nil
}
func (m *mockAuthServiceClient) UpdateUserProfile(ctx context.Context, in *authnservicev1.UpdateUserProfileRequest, opts ...grpc.CallOption) (*authnservicev1.UpdateUserProfileResponse, error) {
	return &authnservicev1.UpdateUserProfileResponse{}, nil
}
func (m *mockAuthServiceClient) UploadUserDocument(ctx context.Context, in *authnservicev1.UploadUserDocumentRequest, opts ...grpc.CallOption) (*authnservicev1.UploadUserDocumentResponse, error) {
	return &authnservicev1.UploadUserDocumentResponse{}, nil
}
func (m *mockAuthServiceClient) GetUserDocument(ctx context.Context, in *authnservicev1.GetUserDocumentRequest, opts ...grpc.CallOption) (*authnservicev1.GetUserDocumentResponse, error) {
	return &authnservicev1.GetUserDocumentResponse{}, nil
}
func (m *mockAuthServiceClient) UpdateUserDocument(ctx context.Context, in *authnservicev1.UpdateUserDocumentRequest, opts ...grpc.CallOption) (*authnservicev1.UpdateUserDocumentResponse, error) {
	return &authnservicev1.UpdateUserDocumentResponse{}, nil
}
func (m *mockAuthServiceClient) ListUserDocuments(ctx context.Context, in *authnservicev1.ListUserDocumentsRequest, opts ...grpc.CallOption) (*authnservicev1.ListUserDocumentsResponse, error) {
	return &authnservicev1.ListUserDocumentsResponse{}, nil
}
func (m *mockAuthServiceClient) DeleteUserDocument(ctx context.Context, in *authnservicev1.DeleteUserDocumentRequest, opts ...grpc.CallOption) (*authnservicev1.DeleteUserDocumentResponse, error) {
	return &authnservicev1.DeleteUserDocumentResponse{}, nil
}
func (m *mockAuthServiceClient) ListDocumentTypes(ctx context.Context, in *authnservicev1.ListDocumentTypesRequest, opts ...grpc.CallOption) (*authnservicev1.ListDocumentTypesResponse, error) {
	return &authnservicev1.ListDocumentTypesResponse{}, nil
}
func (m *mockAuthServiceClient) InitiateKYC(ctx context.Context, in *authnservicev1.InitiateKYCRequest, opts ...grpc.CallOption) (*authnservicev1.InitiateKYCResponse, error) {
	return &authnservicev1.InitiateKYCResponse{}, nil
}
func (m *mockAuthServiceClient) GetKYCStatus(ctx context.Context, in *authnservicev1.GetKYCStatusRequest, opts ...grpc.CallOption) (*authnservicev1.GetKYCStatusResponse, error) {
	return &authnservicev1.GetKYCStatusResponse{}, nil
}
func (m *mockAuthServiceClient) SubmitKYCFrame(ctx context.Context, in *authnservicev1.SubmitKYCFrameRequest, opts ...grpc.CallOption) (*authnservicev1.SubmitKYCFrameResponse, error) {
	return &authnservicev1.SubmitKYCFrameResponse{}, nil
}
func (m *mockAuthServiceClient) ApproveKYC(ctx context.Context, in *authnservicev1.ApproveKYCRequest, opts ...grpc.CallOption) (*authnservicev1.ApproveKYCResponse, error) {
	return &authnservicev1.ApproveKYCResponse{}, nil
}
func (m *mockAuthServiceClient) RejectKYC(ctx context.Context, in *authnservicev1.RejectKYCRequest, opts ...grpc.CallOption) (*authnservicev1.RejectKYCResponse, error) {
	return &authnservicev1.RejectKYCResponse{}, nil
}
func (m *mockAuthServiceClient) CompleteKYCSession(ctx context.Context, in *authnservicev1.CompleteKYCSessionRequest, opts ...grpc.CallOption) (*authnservicev1.CompleteKYCSessionResponse, error) {
	return &authnservicev1.CompleteKYCSessionResponse{}, nil
}
func (m *mockAuthServiceClient) VerifyDocument(ctx context.Context, in *authnservicev1.VerifyDocumentRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyDocumentResponse, error) {
	return &authnservicev1.VerifyDocumentResponse{}, nil
}
func (m *mockAuthServiceClient) CreateVoiceSession(ctx context.Context, in *authnservicev1.CreateVoiceSessionRequest, opts ...grpc.CallOption) (*authnservicev1.CreateVoiceSessionResponse, error) {
	return &authnservicev1.CreateVoiceSessionResponse{}, nil
}
func (m *mockAuthServiceClient) GetVoiceSession(ctx context.Context, in *authnservicev1.GetVoiceSessionRequest, opts ...grpc.CallOption) (*authnservicev1.GetVoiceSessionResponse, error) {
	return &authnservicev1.GetVoiceSessionResponse{}, nil
}
func (m *mockAuthServiceClient) EndVoiceSession(ctx context.Context, in *authnservicev1.EndVoiceSessionRequest, opts ...grpc.CallOption) (*authnservicev1.EndVoiceSessionResponse, error) {
	return &authnservicev1.EndVoiceSessionResponse{}, nil
}
func (m *mockAuthServiceClient) GetProfilePhotoUploadURL(ctx context.Context, in *authnservicev1.GetProfilePhotoUploadURLRequest, opts ...grpc.CallOption) (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
	return &authnservicev1.GetProfilePhotoUploadURLResponse{}, nil
}
func (m *mockAuthServiceClient) UpdateNotificationPreferences(ctx context.Context, in *authnservicev1.UpdateNotificationPreferencesRequest, opts ...grpc.CallOption) (*authnservicev1.UpdateNotificationPreferencesResponse, error) {
	return &authnservicev1.UpdateNotificationPreferencesResponse{}, nil
}
func (m *mockAuthServiceClient) EnableTOTP(ctx context.Context, in *authnservicev1.EnableTOTPRequest, opts ...grpc.CallOption) (*authnservicev1.EnableTOTPResponse, error) {
	return &authnservicev1.EnableTOTPResponse{}, nil
}
func (m *mockAuthServiceClient) VerifyTOTP(ctx context.Context, in *authnservicev1.VerifyTOTPRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyTOTPResponse, error) {
	return &authnservicev1.VerifyTOTPResponse{}, nil
}
func (m *mockAuthServiceClient) DisableTOTP(ctx context.Context, in *authnservicev1.DisableTOTPRequest, opts ...grpc.CallOption) (*authnservicev1.DisableTOTPResponse, error) {
	return &authnservicev1.DisableTOTPResponse{}, nil
}
func (m *mockAuthServiceClient) GetJWKS(ctx context.Context, in *authnservicev1.GetJWKSRequest, opts ...grpc.CallOption) (*authnservicev1.GetJWKSResponse, error) {
	return &authnservicev1.GetJWKSResponse{}, nil
}
func (m *mockAuthServiceClient) InitiateVoiceSession(ctx context.Context, in *authnservicev1.InitiateVoiceSessionRequest, opts ...grpc.CallOption) (*authnservicev1.InitiateVoiceSessionResponse, error) {
	return &authnservicev1.InitiateVoiceSessionResponse{}, nil
}
func (m *mockAuthServiceClient) SubmitVoiceSample(ctx context.Context, in *authnservicev1.SubmitVoiceSampleRequest, opts ...grpc.CallOption) (*authnservicev1.SubmitVoiceSampleResponse, error) {
	return &authnservicev1.SubmitVoiceSampleResponse{}, nil
}
func (m *mockAuthServiceClient) VerifyVoiceSession(ctx context.Context, in *authnservicev1.VerifyVoiceSessionRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyVoiceSessionResponse, error) {
	return &authnservicev1.VerifyVoiceSessionResponse{}, nil
}

func TestAuthnHandler_Login_SetsSessionCookieForServerSide(t *testing.T) {
	m := &mockAuthServiceClient{}
	m.loginFn = func(ctx context.Context, in *authnservicev1.LoginRequest, _ ...grpc.CallOption) (*authnservicev1.LoginResponse, error) {
		return &authnservicev1.LoginResponse{
			UserId:       "u1",
			SessionId:    "s1",
			SessionType:  "SERVER_SIDE",
			SessionToken: "secret",
			CsrfToken:    "csrf",
		}, nil
	}
	// Inject the mock client directly.
	h := &AuthnHandler{client: m}

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"mobile_number":"+8801","password":"x"}`))
	w := httptest.NewRecorder()

	h.Login(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	cookies := res.Cookies()
	require.NotEmpty(t, cookies)
	require.Equal(t, sessionCookieName, cookies[0].Name)
	require.Equal(t, "secret", cookies[0].Value)
	require.True(t, cookies[0].HttpOnly)
	require.Equal(t, "csrf", res.Header.Get("X-CSRF-Token"))
}

// compile-time guard: ensure handler can work with any proto message
var _ = proto.Marshal
