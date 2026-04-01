package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// AuthService handles auth-related API calls
type AuthService struct {
	Client Client
}

// ListAPIKeys List API keys for an owner
func (s *AuthService) ListAPIKeys(ctx context.Context) error {
	path := "/v1/auth/api-keys"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateAPIKey Create a new API key for a user or service
func (s *AuthService) CreateAPIKey(ctx context.Context, req *models.APIKeyCreationRequest) error {
	path := "/v1/auth/api-keys"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RevokeAPIKey Revoke an API key
func (s *AuthService) RevokeAPIKey(ctx context.Context, keyId string, req *models.RevokeAPIKeyRequest) error {
	path := "/v1/auth/api-keys/{key_id}:revoke"
	path = strings.ReplaceAll(path, "{key_id}", keyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RotateAPIKey Rotate an API key (generates new key, marks old one for graceful expiry)
func (s *AuthService) RotateAPIKey(ctx context.Context, keyId string, req *models.APIKeyRotationRequest) error {
	path := "/v1/auth/api-keys/{key_id}:rotate"
	path = strings.ReplaceAll(path, "{key_id}", keyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// BiometricAuthenticate Authenticate using a device-bound biometric token (mobile only)
func (s *AuthService) BiometricAuthenticate(ctx context.Context, req *models.BiometricAuthenticateRequest) error {
	path := "/v1/auth/biometric:authenticate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ValidateCSRF Validate CSRF token (server-side sessions only)
func (s *AuthService) ValidateCSRF(ctx context.Context, req *models.CSRFValidationRequest) error {
	path := "/v1/auth/csrf:validate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListDocumentTypes List document types
func (s *AuthService) ListDocumentTypes(ctx context.Context) error {
	path := "/v1/auth/document-types"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetUserDocument Get user document
func (s *AuthService) GetUserDocument(ctx context.Context, userDocumentId string) error {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateUserDocument Update user document
func (s *AuthService) UpdateUserDocument(ctx context.Context, userDocumentId string, req *models.UserDocumentUpdateRequest) error {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeleteUserDocument Delete user document
func (s *AuthService) DeleteUserDocument(ctx context.Context, userDocumentId string) error {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// VerifyDocument ── Document Verification (Admin) ──
func (s *AuthService) VerifyDocument(ctx context.Context, userDocumentId string, req *models.DocumentVerificationRequest) error {
	path := "/v1/auth/documents/{user_document_id}:verify"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// EmailPasswordLogin Login via email + password (B2B beneficiary self-service)
func (s *AuthService) EmailPasswordLogin(ctx context.Context, req *models.EmailPasswordLoginRequest) error {
	path := "/v1/auth/email-password:login"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// EmailLogin Login via email + OTP (Business Beneficiary / System User only →
func (s *AuthService) EmailLogin(ctx context.Context, req *models.EmailLoginRequest) error {
	path := "/v1/auth/email/login"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendEmailOTP Send email OTP (verification or login)
func (s *AuthService) SendEmailOTP(ctx context.Context, req *models.EmailOTPSendingRequest) error {
	path := "/v1/auth/email/otp:send"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ResetPasswordByEmail Complete password reset using email OTP
func (s *AuthService) ResetPasswordByEmail(ctx context.Context, req *models.ResetPasswordByEmailRequest) error {
	path := "/v1/auth/email/password:reset"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RequestPasswordResetByEmail Request password reset via email OTP
func (s *AuthService) RequestPasswordResetByEmail(ctx context.Context, req *models.RequestPasswordResetByEmailRequest) error {
	path := "/v1/auth/email/password:reset-request"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RegisterEmailUser Register a portal user with email (requires email, triggers email
func (s *AuthService) RegisterEmailUser(ctx context.Context, req *models.EmailUserRegistrationRequest) error {
	path := "/v1/auth/email/register"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyEmail Verify email address using OTP (must call before email login is allowed)
func (s *AuthService) VerifyEmail(ctx context.Context, req *models.EmailVerificationRequest) error {
	path := "/v1/auth/email/verify"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ApproveKYC Approve KYC
func (s *AuthService) ApproveKYC(ctx context.Context, kycId string, req *models.KYCApprovalRequest) error {
	path := "/v1/auth/kyc/{kyc_id}:approve"
	path = strings.ReplaceAll(path, "{kyc_id}", kycId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// Login Login with credentials
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) error {
	path := "/v1/auth/login"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// Logout Logout
func (s *AuthService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	path := "/v1/auth/logout"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ResendOTP Resend OTP (invalidates previous OTP, generates fresh one)
func (s *AuthService) ResendOTP(ctx context.Context, req *models.ResendOTPRequest) error {
	path := "/v1/auth/otp:resend"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SendOTP Send OTP for verification
func (s *AuthService) SendOTP(ctx context.Context, req *models.OTPSendingRequest) error {
	path := "/v1/auth/otp:send"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyOTP Verify OTP
func (s *AuthService) VerifyOTP(ctx context.Context, req *models.OTPVerificationRequest) error {
	path := "/v1/auth/otp:verify"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ChangePassword Change password
func (s *AuthService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	path := "/v1/auth/password:change"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ResetPassword Reset password
func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	path := "/v1/auth/password:reset"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// Register Register new user
func (s *AuthService) Register(ctx context.Context, req *models.RegistrationRequest) error {
	path := "/v1/auth/register"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetCurrentSession Get current user's active session
func (s *AuthService) GetCurrentSession(ctx context.Context) error {
	path := "/v1/auth/session/current"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetSession Get session details
func (s *AuthService) GetSession(ctx context.Context, sessionId string) error {
	path := "/v1/auth/sessions/{session_id}"
	path = strings.ReplaceAll(path, "{session_id}", sessionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RevokeSession Revoke a specific session
func (s *AuthService) RevokeSession(ctx context.Context, sessionId string) error {
	path := "/v1/auth/sessions/{session_id}"
	path = strings.ReplaceAll(path, "{session_id}", sessionId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// RefreshToken Refresh access token
func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) error {
	path := "/v1/auth/token:refresh"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ValidateToken Validate token
func (s *AuthService) ValidateToken(ctx context.Context, req *models.TokenValidationRequest) error {
	path := "/v1/auth/token:validate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListUserDocuments List user documents
func (s *AuthService) ListUserDocuments(ctx context.Context, userId string) error {
	path := "/v1/auth/users/{user_id}/documents"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UploadUserDocument Upload user document
func (s *AuthService) UploadUserDocument(ctx context.Context, userId string, req *models.UserDocumentUploadRequest) error {
	path := "/v1/auth/users/{user_id}/documents"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetKYCStatus Get KYC status
func (s *AuthService) GetKYCStatus(ctx context.Context, userId string) error {
	path := "/v1/auth/users/{user_id}/kyc"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// InitiateKYC ── KYC Verification ──
func (s *AuthService) InitiateKYC(ctx context.Context, userId string, req *models.InitiateKYCRequest) error {
	path := "/v1/auth/users/{user_id}/kyc"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CompleteKYCSession Complete KYC session
func (s *AuthService) CompleteKYCSession(ctx context.Context, userId string, req *models.KYCSessionCompletionRequest) error {
	path := "/v1/auth/users/{user_id}/kyc:complete"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SubmitKYCFrame Submit KYC frame
func (s *AuthService) SubmitKYCFrame(ctx context.Context, userId string, req *models.KYCFrameSubmissionRequest) error {
	path := "/v1/auth/users/{user_id}/kyc:submit-frame"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetNotificationPreferences ── Notification Preferences ──
func (s *AuthService) GetNotificationPreferences(ctx context.Context, userId string) error {
	path := "/v1/auth/users/{user_id}/notification-preferences"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateNotificationPreferences Update notification preferences
func (s *AuthService) UpdateNotificationPreferences(ctx context.Context, userId string, req *models.NotificationPreferencesUpdateRequest) error {
	path := "/v1/auth/users/{user_id}/notification-preferences"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// SetTemporaryPassword Set a temporary password and require the user to change it on next login
func (s *AuthService) SetTemporaryPassword(ctx context.Context, userId string, req *models.SetTemporaryPasswordRequest) error {
	path := "/v1/auth/users/{user_id}/password:temporary"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetUserProfile Get user profile
func (s *AuthService) GetUserProfile(ctx context.Context, userId string) error {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateUserProfile Create user profile
func (s *AuthService) CreateUserProfile(ctx context.Context, userId string, req *models.UserProfileCreationRequest) error {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateUserProfile Update user profile
func (s *AuthService) UpdateUserProfile(ctx context.Context, userId string, req *models.UserProfileUpdateRequest) error {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// GetProfilePhotoUploadURL ── Profile Photo Upload URL ──
func (s *AuthService) GetProfilePhotoUploadURL(ctx context.Context, userId string, req *models.ProfilePhotoUploadURLRetrievalRequest) error {
	path := "/v1/auth/users/{user_id}/profile/photo:upload-url"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListSessions List all sessions for a user
func (s *AuthService) ListSessions(ctx context.Context, userId string) error {
	path := "/v1/auth/users/{user_id}/sessions"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RevokeAllSessions Revoke all sessions for a user (logout from all devices)
func (s *AuthService) RevokeAllSessions(ctx context.Context, userId string, req *models.RevokeAllSessionsRequest) error {
	path := "/v1/auth/users/{user_id}/sessions:revoke-all"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DisableTOTP Disable TOTP
func (s *AuthService) DisableTOTP(ctx context.Context, userId string, req *models.TOTPDisablementRequest) error {
	path := "/v1/auth/users/{user_id}/totp:disable"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// EnableTOTP 🔐 TOTP / 2FA 🔐
func (s *AuthService) EnableTOTP(ctx context.Context, userId string, req *models.TOTPEnablementRequest) error {
	path := "/v1/auth/users/{user_id}/totp:enable"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyTOTP Verify TOTP
func (s *AuthService) VerifyTOTP(ctx context.Context, userId string, req *models.TOTPVerificationRequest) error {
	path := "/v1/auth/users/{user_id}/totp:verify"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// FindPortalUser Find a portal user by exact email or mobile number
func (s *AuthService) FindPortalUser(ctx context.Context, req *models.FindPortalUserRequest) error {
	path := "/v1/auth/users:find"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// InitiateVoiceSession 🎤 Voice Biometric Auth (Sprint 1
func (s *AuthService) InitiateVoiceSession(ctx context.Context, req *models.InitiateVoiceSessionRequest) error {
	path := "/v1/auth/voice-biometric:initiate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SubmitVoiceSample Submit voice sample
func (s *AuthService) SubmitVoiceSample(ctx context.Context, req *models.VoiceSampleSubmissionRequest) error {
	path := "/v1/auth/voice-biometric:submit"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyVoiceSession Verify voice session
func (s *AuthService) VerifyVoiceSession(ctx context.Context, req *models.VoiceSessionVerificationRequest) error {
	path := "/v1/auth/voice-biometric:verify"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CreateVoiceSession ── Voice Sessions ──
func (s *AuthService) CreateVoiceSession(ctx context.Context, req *models.VoiceSessionCreationRequest) error {
	path := "/v1/auth/voice-sessions"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

