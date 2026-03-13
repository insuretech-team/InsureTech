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

// SendOTP Send OTP for verification
func (s *AuthService) SendOTP(ctx context.Context, req *models.OTPSendingRequest) (*models.OTPSendingResponse, error) {
	path := "/v1/auth/otp:send"
	var result models.OTPSendingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Logout Logout
func (s *AuthService) Logout(ctx context.Context, req *models.LogoutRequest) (*models.LogoutResponse, error) {
	path := "/v1/auth/logout"
	var result models.LogoutResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyEmail Verify email address using OTP (must call before email login is allowed)
func (s *AuthService) VerifyEmail(ctx context.Context, req *models.EmailVerificationRequest) (*models.EmailVerificationResponse, error) {
	path := "/v1/auth/email/verify"
	var result models.EmailVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadUserDocument Upload user document
func (s *AuthService) UploadUserDocument(ctx context.Context, userId string, req *models.UserDocumentUploadRequest) (*models.UserDocumentUploadResponse, error) {
	path := "/v1/auth/users/{user_id}/documents"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserDocumentUploadResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserDocuments List user documents
func (s *AuthService) ListUserDocuments(ctx context.Context, userId string) (*models.UserDocumentsListingResponse, error) {
	path := "/v1/auth/users/{user_id}/documents"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserDocumentsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKYCStatus Get k y c status
func (s *AuthService) GetKYCStatus(ctx context.Context, userId string) (*models.KYCStatusRetrievalResponse, error) {
	path := "/v1/auth/users/{user_id}/kyc"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.KYCStatusRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiateKYC ── KYC Verification ──
func (s *AuthService) InitiateKYC(ctx context.Context, userId string, req *models.InitiateKYCRequest) (*models.InitiateKYCResponse, error) {
	path := "/v1/auth/users/{user_id}/kyc"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.InitiateKYCResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CompleteKYCSession Complete k y c session
func (s *AuthService) CompleteKYCSession(ctx context.Context, userId string, req *models.KYCSessionCompletionRequest) (*models.KYCSessionCompletionResponse, error) {
	path := "/v1/auth/users/{user_id}/kyc:complete"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.KYCSessionCompletionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EndVoiceSession End voice session
func (s *AuthService) EndVoiceSession(ctx context.Context, voiceSessionId string, req *models.AuthnEndVoiceSessionRequest) (*models.AuthnEndVoiceSessionResponse, error) {
	path := "/v1/auth/voice-sessions/{voice_session_id}:end"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.AuthnEndVoiceSessionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiateVoiceSession 🎤 Voice Biometric Auth (Sprint 1
func (s *AuthService) InitiateVoiceSession(ctx context.Context, req *models.InitiateVoiceSessionRequest) (*models.InitiateVoiceSessionResponse, error) {
	path := "/v1/auth/voice-biometric:initiate"
	var result models.InitiateVoiceSessionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetJWKS 🔑 JWKS 🔑
func (s *AuthService) GetJWKS(ctx context.Context) (*models.AuthnJWKSRetrievalResponse, error) {
	path := "/v1/auth/.well-known/jwks.json"
	var result models.AuthnJWKSRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Register Register new user
func (s *AuthService) Register(ctx context.Context, req *models.RegistrationRequest) (*models.RegistrationResponse, error) {
	path := "/v1/auth/register"
	var result models.RegistrationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyVoiceSession Verify voice session
func (s *AuthService) VerifyVoiceSession(ctx context.Context, req *models.VoiceSessionVerificationRequest) (*models.VoiceSessionVerificationResponse, error) {
	path := "/v1/auth/voice-biometric:verify"
	var result models.VoiceSessionVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValidateCSRF Validate CSRF token (server-side sessions only)
func (s *AuthService) ValidateCSRF(ctx context.Context, req *models.CSRFValidationRequest) (*models.CSRFValidationResponse, error) {
	path := "/v1/auth/csrf:validate"
	var result models.CSRFValidationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SendEmailOTP Send email OTP (verification or login)
func (s *AuthService) SendEmailOTP(ctx context.Context, req *models.EmailOTPSendingRequest) (*models.EmailOTPSendingResponse, error) {
	path := "/v1/auth/email/otp:send"
	var result models.EmailOTPSendingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestPasswordResetByEmail Request password reset via email OTP
func (s *AuthService) RequestPasswordResetByEmail(ctx context.Context, req *models.RequestPasswordResetByEmailRequest) (*models.RequestPasswordResetByEmailResponse, error) {
	path := "/v1/auth/email/password:reset-request"
	var result models.RequestPasswordResetByEmailResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ApproveKYC Approve k y c
func (s *AuthService) ApproveKYC(ctx context.Context, kycId string, req *models.KYCApprovalRequest) (*models.KYCApprovalResponse, error) {
	path := "/v1/auth/kyc/{kyc_id}:approve"
	path = strings.ReplaceAll(path, "{kyc_id}", kycId)
	var result models.KYCApprovalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListSessions List all sessions for a user
func (s *AuthService) ListSessions(ctx context.Context, userId string) (*models.SessionsListingResponse, error) {
	path := "/v1/auth/users/{user_id}/sessions"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.SessionsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RefreshToken Refresh access token
func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	path := "/v1/auth/token:refresh"
	var result models.RefreshTokenResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeSession Revoke a specific session
func (s *AuthService) RevokeSession(ctx context.Context, sessionId string) error {
	path := "/v1/auth/sessions/{session_id}"
	path = strings.ReplaceAll(path, "{session_id}", sessionId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetSession Get session details
func (s *AuthService) GetSession(ctx context.Context, sessionId string) (*models.SessionRetrievalResponse, error) {
	path := "/v1/auth/sessions/{session_id}"
	path = strings.ReplaceAll(path, "{session_id}", sessionId)
	var result models.SessionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyTOTP Verify t o t p
func (s *AuthService) VerifyTOTP(ctx context.Context, userId string, req *models.TOTPVerificationRequest) (*models.TOTPVerificationResponse, error) {
	path := "/v1/auth/users/{user_id}/totp:verify"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.TOTPVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ResetPassword Reset password
func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	path := "/v1/auth/password:reset"
	var result models.ResetPasswordResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ResetPasswordByEmail Complete password reset using email OTP
func (s *AuthService) ResetPasswordByEmail(ctx context.Context, req *models.ResetPasswordByEmailRequest) (*models.ResetPasswordByEmailResponse, error) {
	path := "/v1/auth/email/password:reset"
	var result models.ResetPasswordByEmailResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocumentTypes List document types
func (s *AuthService) ListDocumentTypes(ctx context.Context) (*models.DocumentTypesListingResponse, error) {
	path := "/v1/auth/document-types"
	var result models.DocumentTypesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAPIKey Create a new API key for a user or service
func (s *AuthService) CreateAPIKey(ctx context.Context, req *models.APIKeyCreationRequest) (*models.APIKeyCreationResponse, error) {
	path := "/v1/auth/api-keys"
	var result models.APIKeyCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAPIKeys List API keys for an owner
func (s *AuthService) ListAPIKeys(ctx context.Context) (*models.APIKeysListingResponse, error) {
	path := "/v1/auth/api-keys"
	var result models.APIKeysListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCurrentSession Get current user's active session
func (s *AuthService) GetCurrentSession(ctx context.Context) (*models.CurrentSessionRetrievalResponse, error) {
	path := "/v1/auth/session/current"
	var result models.CurrentSessionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyOTP Verify OTP
func (s *AuthService) VerifyOTP(ctx context.Context, req *models.OTPVerificationRequest) (*models.OTPVerificationResponse, error) {
	path := "/v1/auth/otp:verify"
	var result models.OTPVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SubmitKYCFrame Submit k y c frame
func (s *AuthService) SubmitKYCFrame(ctx context.Context, userId string, req *models.KYCFrameSubmissionRequest) (*models.KYCFrameSubmissionResponse, error) {
	path := "/v1/auth/users/{user_id}/kyc:submit-frame"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.KYCFrameSubmissionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeAllSessions Revoke all sessions for a user (logout from all devices)
func (s *AuthService) RevokeAllSessions(ctx context.Context, userId string, req *models.RevokeAllSessionsRequest) (*models.RevokeAllSessionsResponse, error) {
	path := "/v1/auth/users/{user_id}/sessions:revoke-all"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.RevokeAllSessionsResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserDocument Get user document
func (s *AuthService) GetUserDocument(ctx context.Context, userDocumentId string) (*models.UserDocumentRetrievalResponse, error) {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	var result models.UserDocumentRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateUserDocument Update user document
func (s *AuthService) UpdateUserDocument(ctx context.Context, userDocumentId string, req *models.UserDocumentUpdateRequest) (*models.UserDocumentUpdateResponse, error) {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	var result models.UserDocumentUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteUserDocument Delete user document
func (s *AuthService) DeleteUserDocument(ctx context.Context, userDocumentId string) error {
	path := "/v1/auth/documents/{user_document_id}"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// RotateAPIKey Rotate an API key (generates new key, marks old one for graceful expiry)
func (s *AuthService) RotateAPIKey(ctx context.Context, keyId string, req *models.APIKeyRotationRequest) (*models.APIKeyRotationResponse, error) {
	path := "/v1/auth/api-keys/{key_id}:rotate"
	path = strings.ReplaceAll(path, "{key_id}", keyId)
	var result models.APIKeyRotationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ChangePassword Change password
func (s *AuthService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	path := "/v1/auth/password:change"
	var result models.ChangePasswordResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVoiceSession Get voice session
func (s *AuthService) GetVoiceSession(ctx context.Context, voiceSessionId string) (*models.AuthnVoiceSessionRetrievalResponse, error) {
	path := "/v1/auth/voice-sessions/{voice_session_id}"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.AuthnVoiceSessionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SubmitVoiceSample Submit voice sample
func (s *AuthService) SubmitVoiceSample(ctx context.Context, req *models.VoiceSampleSubmissionRequest) (*models.VoiceSampleSubmissionResponse, error) {
	path := "/v1/auth/voice-biometric:submit"
	var result models.VoiceSampleSubmissionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterEmailUser Register a portal user with email (requires email, triggers email
func (s *AuthService) RegisterEmailUser(ctx context.Context, req *models.EmailUserRegistrationRequest) (*models.EmailUserRegistrationResponse, error) {
	path := "/v1/auth/email/register"
	var result models.EmailUserRegistrationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ResendOTP Resend OTP (invalidates previous OTP, generates fresh one)
func (s *AuthService) ResendOTP(ctx context.Context, req *models.ResendOTPRequest) (*models.ResendOTPResponse, error) {
	path := "/v1/auth/otp:resend"
	var result models.ResendOTPResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateUserProfile Create user profile
func (s *AuthService) CreateUserProfile(ctx context.Context, userId string, req *models.UserProfileCreationRequest) (*models.UserProfileCreationResponse, error) {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserProfileCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserProfile Get user profile
func (s *AuthService) GetUserProfile(ctx context.Context, userId string) (*models.UserProfileRetrievalResponse, error) {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserProfileRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateUserProfile Update user profile
func (s *AuthService) UpdateUserProfile(ctx context.Context, userId string, req *models.UserProfileUpdateRequest) (*models.UserProfileUpdateResponse, error) {
	path := "/v1/auth/users/{user_id}/profile"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.UserProfileUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DisableTOTP Disable t o t p
func (s *AuthService) DisableTOTP(ctx context.Context, userId string, req *models.TOTPDisablementRequest) (*models.TOTPDisablementResponse, error) {
	path := "/v1/auth/users/{user_id}/totp:disable"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.TOTPDisablementResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProfilePhotoUploadURL ── Profile Photo Upload URL ──
func (s *AuthService) GetProfilePhotoUploadURL(ctx context.Context, userId string, req *models.ProfilePhotoUploadURLRetrievalRequest) (*models.ProfilePhotoUploadURLRetrievalResponse, error) {
	path := "/v1/auth/users/{user_id}/profile/photo:upload-url"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.ProfilePhotoUploadURLRetrievalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyDocument ── Document Verification (Admin) ──
func (s *AuthService) VerifyDocument(ctx context.Context, userDocumentId string, req *models.DocumentVerificationRequest) (*models.DocumentVerificationResponse, error) {
	path := "/v1/auth/documents/{user_document_id}:verify"
	path = strings.ReplaceAll(path, "{user_document_id}", userDocumentId)
	var result models.DocumentVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Login Login with credentials
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	path := "/v1/auth/login"
	var result models.LoginResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BiometricAuthenticate Authenticate using a device-bound biometric token (mobile only)
func (s *AuthService) BiometricAuthenticate(ctx context.Context, req *models.BiometricAuthenticateRequest) (*models.BiometricAuthenticateResponse, error) {
	path := "/v1/auth/biometric:authenticate"
	var result models.BiometricAuthenticateResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateVoiceSession ── Voice Sessions ──
func (s *AuthService) CreateVoiceSession(ctx context.Context, req *models.VoiceSessionCreationRequest) (*models.VoiceSessionCreationResponse, error) {
	path := "/v1/auth/voice-sessions"
	var result models.VoiceSessionCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EnableTOTP 🔐 TOTP / 2FA 🔐
func (s *AuthService) EnableTOTP(ctx context.Context, userId string, req *models.TOTPEnablementRequest) (*models.TOTPEnablementResponse, error) {
	path := "/v1/auth/users/{user_id}/totp:enable"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.TOTPEnablementResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValidateToken Validate token
func (s *AuthService) ValidateToken(ctx context.Context, req *models.TokenValidationRequest) (*models.TokenValidationResponse, error) {
	path := "/v1/auth/token:validate"
	var result models.TokenValidationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateNotificationPreferences ── Notification Preferences ──
func (s *AuthService) UpdateNotificationPreferences(ctx context.Context, userId string, req *models.NotificationPreferencesUpdateRequest) (*models.NotificationPreferencesUpdateResponse, error) {
	path := "/v1/auth/users/{user_id}/notification-preferences"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.NotificationPreferencesUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EmailLogin Login via email + OTP (Business Beneficiary / System User only →
func (s *AuthService) EmailLogin(ctx context.Context, req *models.EmailLoginRequest) (*models.EmailLoginResponse, error) {
	path := "/v1/auth/email/login"
	var result models.EmailLoginResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeAPIKey Revoke an API key
func (s *AuthService) RevokeAPIKey(ctx context.Context, keyId string, req *models.RevokeAPIKeyRequest) (*models.RevokeAPIKeyResponse, error) {
	path := "/v1/auth/api-keys/{key_id}:revoke"
	path = strings.ReplaceAll(path, "{key_id}", keyId)
	var result models.RevokeAPIKeyResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RejectKYC Reject k y c
func (s *AuthService) RejectKYC(ctx context.Context, kycId string, req *models.AuthnKYCRejectionRequest) (*models.AuthnKYCRejectionResponse, error) {
	path := "/v1/auth/kyc/{kyc_id}:reject"
	path = strings.ReplaceAll(path, "{kyc_id}", kycId)
	var result models.AuthnKYCRejectionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

