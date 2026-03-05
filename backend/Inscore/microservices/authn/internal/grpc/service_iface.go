package grpc

import (
	"context"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
)

// AuthServiceIface is the interface that AuthServiceHandler depends on.
// Using an interface instead of a concrete *service.AuthService allows the handler
// to be tested without a real database or Redis connection.
type AuthServiceIface interface {
	Register(ctx context.Context, req *authnservicev1.RegisterRequest) (*authnservicev1.RegisterResponse, error)
	SendOTP(ctx context.Context, req *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error)
	VerifyOTP(ctx context.Context, req *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error)
	ResendOTP(ctx context.Context, req *authnservicev1.ResendOTPRequest) (*authnservicev1.ResendOTPResponse, error)
	Login(ctx context.Context, req *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error)
	RefreshToken(ctx context.Context, req *authnservicev1.RefreshTokenRequest) (*authnservicev1.RefreshTokenResponse, error)
	Logout(ctx context.Context, req *authnservicev1.LogoutRequest) (*authnservicev1.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *authnservicev1.ChangePasswordRequest) (*authnservicev1.ChangePasswordResponse, error)
	ResetPassword(ctx context.Context, req *authnservicev1.ResetPasswordRequest) (*authnservicev1.ResetPasswordResponse, error)
	ValidateToken(ctx context.Context, req *authnservicev1.ValidateTokenRequest) (*authnservicev1.ValidateTokenResponse, error)
	GetSession(ctx context.Context, req *authnservicev1.GetSessionRequest) (*authnservicev1.GetSessionResponse, error)
	ListSessions(ctx context.Context, req *authnservicev1.ListSessionsRequest) (*authnservicev1.ListSessionsResponse, error)
	RevokeSession(ctx context.Context, req *authnservicev1.RevokeSessionRequest) (*authnservicev1.RevokeSessionResponse, error)
	ValidateCSRF(ctx context.Context, req *authnservicev1.ValidateCSRFRequest) (*authnservicev1.ValidateCSRFResponse, error)
	GetCurrentSession(ctx context.Context, req *authnservicev1.GetCurrentSessionRequest) (*authnservicev1.GetCurrentSessionResponse, error)
	RevokeAllSessions(ctx context.Context, req *authnservicev1.RevokeAllSessionsRequest) (*authnservicev1.RevokeAllSessionsResponse, error)
	RegisterEmailUser(ctx context.Context, req *authnservicev1.RegisterEmailUserRequest) (*authnservicev1.RegisterEmailUserResponse, error)
	SendEmailOTP(ctx context.Context, req *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error)
	VerifyEmail(ctx context.Context, req *authnservicev1.VerifyEmailRequest) (*authnservicev1.VerifyEmailResponse, error)
	EmailLogin(ctx context.Context, req *authnservicev1.EmailLoginRequest) (*authnservicev1.EmailLoginResponse, error)
	RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error)
	ResetPasswordByEmail(ctx context.Context, req *authnservicev1.ResetPasswordByEmailRequest) (*authnservicev1.ResetPasswordByEmailResponse, error)
	BiometricAuthenticate(ctx context.Context, req *authnservicev1.BiometricAuthenticateRequest) (*authnservicev1.BiometricAuthenticateResponse, error)
	UpdateDLRStatus(ctx context.Context, req *authnservicev1.UpdateDLRStatusRequest) (*authnservicev1.UpdateDLRStatusResponse, error)
	CreateAPIKey(ctx context.Context, req *authnservicev1.CreateAPIKeyRequest) (*authnservicev1.CreateAPIKeyResponse, error)
	ListAPIKeys(ctx context.Context, req *authnservicev1.ListAPIKeysRequest) (*authnservicev1.ListAPIKeysResponse, error)
	RevokeAPIKey(ctx context.Context, req *authnservicev1.RevokeAPIKeyRequest) (*authnservicev1.RevokeAPIKeyResponse, error)
	RotateAPIKey(ctx context.Context, req *authnservicev1.RotateAPIKeyRequest) (*authnservicev1.RotateAPIKeyResponse, error)
	// User Profile
	CreateUserProfile(ctx context.Context, req *authnservicev1.CreateUserProfileRequest) (*authnservicev1.CreateUserProfileResponse, error)
	GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, req *authnservicev1.UpdateUserProfileRequest) (*authnservicev1.UpdateUserProfileResponse, error)
	// User Documents
	UploadUserDocument(ctx context.Context, req *authnservicev1.UploadUserDocumentRequest) (*authnservicev1.UploadUserDocumentResponse, error)
	ListUserDocuments(ctx context.Context, req *authnservicev1.ListUserDocumentsRequest) (*authnservicev1.ListUserDocumentsResponse, error)
	GetUserDocument(ctx context.Context, req *authnservicev1.GetUserDocumentRequest) (*authnservicev1.GetUserDocumentResponse, error)
	UpdateUserDocument(ctx context.Context, req *authnservicev1.UpdateUserDocumentRequest) (*authnservicev1.UpdateUserDocumentResponse, error)
	DeleteUserDocument(ctx context.Context, req *authnservicev1.DeleteUserDocumentRequest) (*authnservicev1.DeleteUserDocumentResponse, error)
	// Document Types
	ListDocumentTypes(ctx context.Context, req *authnservicev1.ListDocumentTypesRequest) (*authnservicev1.ListDocumentTypesResponse, error)
	// KYC Verification
	InitiateKYC(ctx context.Context, req *authnservicev1.InitiateKYCRequest) (*authnservicev1.InitiateKYCResponse, error)
	GetKYCStatus(ctx context.Context, req *authnservicev1.GetKYCStatusRequest) (*authnservicev1.GetKYCStatusResponse, error)
	SubmitKYCFrame(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error)
	CompleteKYCSession(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error)
	ApproveKYC(ctx context.Context, req *authnservicev1.ApproveKYCRequest) (*authnservicev1.ApproveKYCResponse, error)
	RejectKYC(ctx context.Context, req *authnservicev1.RejectKYCRequest) (*authnservicev1.RejectKYCResponse, error)
	// Document Verification
	VerifyDocument(ctx context.Context, req *authnservicev1.VerifyDocumentRequest) (*authnservicev1.VerifyDocumentResponse, error)
	// Voice Sessions (CRUD / IVR)
	CreateVoiceSession(ctx context.Context, req *authnservicev1.CreateVoiceSessionRequest) (*authnservicev1.CreateVoiceSessionResponse, error)
	GetVoiceSession(ctx context.Context, req *authnservicev1.GetVoiceSessionRequest) (*authnservicev1.GetVoiceSessionResponse, error)
	EndVoiceSession(ctx context.Context, req *authnservicev1.EndVoiceSessionRequest) (*authnservicev1.EndVoiceSessionResponse, error)
	// Voice Biometric Auth (Sprint 1.10)
	InitiateVoiceSession(ctx context.Context, req *authnservicev1.InitiateVoiceSessionRequest) (*authnservicev1.InitiateVoiceSessionResponse, error)
	SubmitVoiceSample(ctx context.Context, req *authnservicev1.SubmitVoiceSampleRequest) (*authnservicev1.SubmitVoiceSampleResponse, error)
	VerifyVoiceSession(ctx context.Context, req *authnservicev1.VerifyVoiceSessionRequest) (*authnservicev1.VerifyVoiceSessionResponse, error)
	// Profile Photo
	GetProfilePhotoUploadURL(ctx context.Context, req *authnservicev1.GetProfilePhotoUploadURLRequest) (*authnservicev1.GetProfilePhotoUploadURLResponse, error)
	// Notification Preferences
	UpdateNotificationPreferences(ctx context.Context, req *authnservicev1.UpdateNotificationPreferencesRequest) (*authnservicev1.UpdateNotificationPreferencesResponse, error)
	// TOTP / 2FA
	EnableTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error)
	VerifyTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error)
	DisableTOTP(ctx context.Context, req *authnservicev1.DisableTOTPRequest) (*authnservicev1.DisableTOTPResponse, error)
	// JWKS
	GetJWKS(ctx context.Context, req *authnservicev1.GetJWKSRequest) (*authnservicev1.GetJWKSResponse, error)
}
