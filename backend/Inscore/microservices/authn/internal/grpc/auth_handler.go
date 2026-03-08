package grpc

import (
	"context"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// normalizeMobile normalises and validates an inbound mobile number,
// returning the canonical +880XXXXXXXXXX form or a gRPC InvalidArgument error.
func normalizeMobile(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", status.Error(codes.InvalidArgument, "mobile_number is required")
	}
	normalized, err := sms.NormalizePhoneNumber(trimmed)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument,
			"invalid mobile_number %q: must be a valid Bangladesh number "+
				"(e.g. 01712345678, +8801712345678, 008801712345678)", raw)
	}
	return "+" + normalized, nil
}

type AuthServiceHandler struct {
	authnservicev1.UnimplementedAuthServiceServer
	authService AuthServiceIface
}

func NewAuthServiceHandler(authService AuthServiceIface) *AuthServiceHandler {
	return &AuthServiceHandler{
		authService: authService,
	}
}

// ── Phone/OTP flows ──────────────────────────────────────────────────────────

func (h *AuthServiceHandler) Login(ctx context.Context, req *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error) {
	normalized, err := normalizeMobile(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	req.MobileNumber = normalized
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) Register(ctx context.Context, req *authnservicev1.RegisterRequest) (*authnservicev1.RegisterResponse, error) {
	normalized, err := normalizeMobile(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	req.MobileNumber = normalized
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	resp, err := h.authService.Register(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) SendOTP(ctx context.Context, req *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error) {
	if req.Recipient == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	resp, err := h.authService.SendOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) VerifyOTP(ctx context.Context, req *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error) {
	if req.OtpId == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_id is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	resp, err := h.authService.VerifyOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ResendOTP(ctx context.Context, req *authnservicev1.ResendOTPRequest) (*authnservicev1.ResendOTPResponse, error) {
	if req.OriginalOtpId == "" {
		return nil, status.Error(codes.InvalidArgument, "original_otp_id is required")
	}
	resp, err := h.authService.ResendOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Token / Session flows ────────────────────────────────────────────────────

func (h *AuthServiceHandler) ValidateToken(ctx context.Context, req *authnservicev1.ValidateTokenRequest) (*authnservicev1.ValidateTokenResponse, error) {
	if req.AccessToken == "" && req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token or session_id is required")
	}
	resp, err := h.authService.ValidateToken(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RefreshToken(ctx context.Context, req *authnservicev1.RefreshTokenRequest) (*authnservicev1.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}
	resp, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) Logout(ctx context.Context, req *authnservicev1.LogoutRequest) (*authnservicev1.LogoutResponse, error) {
	if req.SessionId == "" && req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id or access_token is required")
	}
	resp, err := h.authService.Logout(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetSession(ctx context.Context, req *authnservicev1.GetSessionRequest) (*authnservicev1.GetSessionResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	resp, err := h.authService.GetSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ListSessions(ctx context.Context, req *authnservicev1.ListSessionsRequest) (*authnservicev1.ListSessionsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.ListSessions(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RevokeSession(ctx context.Context, req *authnservicev1.RevokeSessionRequest) (*authnservicev1.RevokeSessionResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	resp, err := h.authService.RevokeSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RevokeAllSessions(ctx context.Context, req *authnservicev1.RevokeAllSessionsRequest) (*authnservicev1.RevokeAllSessionsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.RevokeAllSessions(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetCurrentSession(ctx context.Context, req *authnservicev1.GetCurrentSessionRequest) (*authnservicev1.GetCurrentSessionResponse, error) {
	// Session token is extracted from auth metadata by the auth interceptor.
	// No additional field validation needed on GetCurrentSession.
	resp, err := h.authService.GetCurrentSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ValidateCSRF(ctx context.Context, req *authnservicev1.ValidateCSRFRequest) (*authnservicev1.ValidateCSRFResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if req.CsrfToken == "" {
		return nil, status.Error(codes.InvalidArgument, "csrf_token is required")
	}
	resp, err := h.authService.ValidateCSRF(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Password flows ───────────────────────────────────────────────────────────

func (h *AuthServiceHandler) ChangePassword(ctx context.Context, req *authnservicev1.ChangePasswordRequest) (*authnservicev1.ChangePasswordResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old_password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	resp, err := h.authService.ChangePassword(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ResetPassword(ctx context.Context, req *authnservicev1.ResetPasswordRequest) (*authnservicev1.ResetPasswordResponse, error) {
	normalized, err := normalizeMobile(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	req.MobileNumber = normalized
	if req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_code is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	resp, err := h.authService.ResetPassword(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Email flows ──────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) RegisterEmailUser(ctx context.Context, req *authnservicev1.RegisterEmailUserRequest) (*authnservicev1.RegisterEmailUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if strings.TrimSpace(req.MobileNumber) != "" {
		normalized, err := normalizeMobile(req.MobileNumber)
		if err != nil {
			return nil, err
		}
		req.MobileNumber = normalized
	}
	resp, err := h.authService.RegisterEmailUser(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) SendEmailOTP(ctx context.Context, req *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	resp, err := h.authService.SendEmailOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) VerifyEmail(ctx context.Context, req *authnservicev1.VerifyEmailRequest) (*authnservicev1.VerifyEmailResponse, error) {
	if req.OtpId == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_id is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	resp, err := h.authService.VerifyEmail(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) EmailLogin(ctx context.Context, req *authnservicev1.EmailLoginRequest) (*authnservicev1.EmailLoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.OtpId == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_id is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	resp, err := h.authService.EmailLogin(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	resp, err := h.authService.RequestPasswordResetByEmail(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ResetPasswordByEmail(ctx context.Context, req *authnservicev1.ResetPasswordByEmailRequest) (*authnservicev1.ResetPasswordByEmailResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.OtpId == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_id is required")
	}
	if req.OtpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "otp_code is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	resp, err := h.authService.ResetPasswordByEmail(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Biometric flow ───────────────────────────────────────────────────────────

func (h *AuthServiceHandler) BiometricAuthenticate(ctx context.Context, req *authnservicev1.BiometricAuthenticateRequest) (*authnservicev1.BiometricAuthenticateResponse, error) {
	if req.BiometricToken == "" {
		return nil, status.Error(codes.InvalidArgument, "biometric_token is required")
	}
	if req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id is required")
	}
	resp, err := h.authService.BiometricAuthenticate(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── DLR webhook ──────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) UpdateDLRStatus(ctx context.Context, req *authnservicev1.UpdateDLRStatusRequest) (*authnservicev1.UpdateDLRStatusResponse, error) {
	if req.ProviderMessageId == "" {
		return nil, status.Error(codes.InvalidArgument, "provider_message_id is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}
	resp, err := h.authService.UpdateDLRStatus(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── API Key management ───────────────────────────────────────────────────────

func (h *AuthServiceHandler) CreateAPIKey(ctx context.Context, req *authnservicev1.CreateAPIKeyRequest) (*authnservicev1.CreateAPIKeyResponse, error) {
	if req.OwnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	resp, err := h.authService.CreateAPIKey(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ListAPIKeys(ctx context.Context, req *authnservicev1.ListAPIKeysRequest) (*authnservicev1.ListAPIKeysResponse, error) {
	if req.OwnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}
	resp, err := h.authService.ListAPIKeys(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RevokeAPIKey(ctx context.Context, req *authnservicev1.RevokeAPIKeyRequest) (*authnservicev1.RevokeAPIKeyResponse, error) {
	if req.KeyId == "" {
		return nil, status.Error(codes.InvalidArgument, "key_id is required")
	}
	resp, err := h.authService.RevokeAPIKey(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RotateAPIKey(ctx context.Context, req *authnservicev1.RotateAPIKeyRequest) (*authnservicev1.RotateAPIKeyResponse, error) {
	if req.KeyId == "" {
		return nil, status.Error(codes.InvalidArgument, "key_id is required")
	}
	resp, err := h.authService.RotateAPIKey(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── User Profile ─────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) CreateUserProfile(ctx context.Context, req *authnservicev1.CreateUserProfileRequest) (*authnservicev1.CreateUserProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.FullName == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name is required")
	}
	if req.AddressLine1 == "" {
		return nil, status.Error(codes.InvalidArgument, "address_line1 is required")
	}
	if req.City == "" {
		return nil, status.Error(codes.InvalidArgument, "city is required")
	}
	if req.District == "" {
		return nil, status.Error(codes.InvalidArgument, "district is required")
	}
	if req.Division == "" {
		return nil, status.Error(codes.InvalidArgument, "division is required")
	}
	if req.NidNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "nid_number is required")
	}
	resp, err := h.authService.CreateUserProfile(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.GetUserProfile(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) UpdateUserProfile(ctx context.Context, req *authnservicev1.UpdateUserProfileRequest) (*authnservicev1.UpdateUserProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.UpdateUserProfile(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── User Documents ───────────────────────────────────────────────────────────

func (h *AuthServiceHandler) UploadUserDocument(ctx context.Context, req *authnservicev1.UploadUserDocumentRequest) (*authnservicev1.UploadUserDocumentResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.DocumentTypeId == "" {
		return nil, status.Error(codes.InvalidArgument, "document_type_id is required")
	}
	if req.FileUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "file_url is required")
	}
	resp, err := h.authService.UploadUserDocument(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ListUserDocuments(ctx context.Context, req *authnservicev1.ListUserDocumentsRequest) (*authnservicev1.ListUserDocumentsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.ListUserDocuments(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetUserDocument(ctx context.Context, req *authnservicev1.GetUserDocumentRequest) (*authnservicev1.GetUserDocumentResponse, error) {
	if req.UserDocumentId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_document_id is required")
	}
	resp, err := h.authService.GetUserDocument(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) UpdateUserDocument(ctx context.Context, req *authnservicev1.UpdateUserDocumentRequest) (*authnservicev1.UpdateUserDocumentResponse, error) {
	if req.UserDocumentId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_document_id is required")
	}
	resp, err := h.authService.UpdateUserDocument(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) DeleteUserDocument(ctx context.Context, req *authnservicev1.DeleteUserDocumentRequest) (*authnservicev1.DeleteUserDocumentResponse, error) {
	if req.UserDocumentId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_document_id is required")
	}
	resp, err := h.authService.DeleteUserDocument(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Document Types ───────────────────────────────────────────────────────────

func (h *AuthServiceHandler) ListDocumentTypes(ctx context.Context, req *authnservicev1.ListDocumentTypesRequest) (*authnservicev1.ListDocumentTypesResponse, error) {
	resp, err := h.authService.ListDocumentTypes(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── KYC Verification ─────────────────────────────────────────────────────────

func (h *AuthServiceHandler) InitiateKYC(ctx context.Context, req *authnservicev1.InitiateKYCRequest) (*authnservicev1.InitiateKYCResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.InitiateKYC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetKYCStatus(ctx context.Context, req *authnservicev1.GetKYCStatusRequest) (*authnservicev1.GetKYCStatusResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.GetKYCStatus(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) SubmitKYCFrame(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if len(req.ImageData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image_data is required")
	}
	resp, err := h.authService.SubmitKYCFrame(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) CompleteKYCSession(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	resp, err := h.authService.CompleteKYCSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) ApproveKYC(ctx context.Context, req *authnservicev1.ApproveKYCRequest) (*authnservicev1.ApproveKYCResponse, error) {
	if req.KycId == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_id is required")
	}
	resp, err := h.authService.ApproveKYC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) RejectKYC(ctx context.Context, req *authnservicev1.RejectKYCRequest) (*authnservicev1.RejectKYCResponse, error) {
	if req.KycId == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_id is required")
	}
	if req.RejectionReason == "" {
		return nil, status.Error(codes.InvalidArgument, "rejection_reason is required")
	}
	resp, err := h.authService.RejectKYC(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Document Verification ────────────────────────────────────────────────────

func (h *AuthServiceHandler) VerifyDocument(ctx context.Context, req *authnservicev1.VerifyDocumentRequest) (*authnservicev1.VerifyDocumentResponse, error) {
	if req.UserDocumentId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_document_id is required")
	}
	resp, err := h.authService.VerifyDocument(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Voice Sessions ───────────────────────────────────────────────────────────

func (h *AuthServiceHandler) CreateVoiceSession(ctx context.Context, req *authnservicev1.CreateVoiceSessionRequest) (*authnservicev1.CreateVoiceSessionResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.CreateVoiceSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) GetVoiceSession(ctx context.Context, req *authnservicev1.GetVoiceSessionRequest) (*authnservicev1.GetVoiceSessionResponse, error) {
	if req.VoiceSessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "voice_session_id is required")
	}
	resp, err := h.authService.GetVoiceSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) EndVoiceSession(ctx context.Context, req *authnservicev1.EndVoiceSessionRequest) (*authnservicev1.EndVoiceSessionResponse, error) {
	if req.VoiceSessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "voice_session_id is required")
	}
	resp, err := h.authService.EndVoiceSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Voice Biometric Auth (Sprint 1.10) ───────────────────────────────────────

func (h *AuthServiceHandler) InitiateVoiceSession(ctx context.Context, req *authnservicev1.InitiateVoiceSessionRequest) (*authnservicev1.InitiateVoiceSessionResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.InitiateVoiceSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) SubmitVoiceSample(ctx context.Context, req *authnservicev1.SubmitVoiceSampleRequest) (*authnservicev1.SubmitVoiceSampleResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	resp, err := h.authService.SubmitVoiceSample(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) VerifyVoiceSession(ctx context.Context, req *authnservicev1.VerifyVoiceSessionRequest) (*authnservicev1.VerifyVoiceSessionResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	resp, err := h.authService.VerifyVoiceSession(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Profile Photo ─────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) GetProfilePhotoUploadURL(ctx context.Context, req *authnservicev1.GetProfilePhotoUploadURLRequest) (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.GetProfilePhotoUploadURL(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── Notification Preferences ──────────────────────────────────────────────────

func (h *AuthServiceHandler) UpdateNotificationPreferences(ctx context.Context, req *authnservicev1.UpdateNotificationPreferencesRequest) (*authnservicev1.UpdateNotificationPreferencesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.UpdateNotificationPreferences(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── TOTP / 2FA ────────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) EnableTOTP(ctx context.Context, req *authnservicev1.EnableTOTPRequest) (*authnservicev1.EnableTOTPResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.EnableTOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) VerifyTOTP(ctx context.Context, req *authnservicev1.VerifyTOTPRequest) (*authnservicev1.VerifyTOTPResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.TotpCode == "" {
		return nil, status.Error(codes.InvalidArgument, "totp_code is required")
	}
	resp, err := h.authService.VerifyTOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

func (h *AuthServiceHandler) DisableTOTP(ctx context.Context, req *authnservicev1.DisableTOTPRequest) (*authnservicev1.DisableTOTPResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	resp, err := h.authService.DisableTOTP(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}

// ── JWKS ─────────────────────────────────────────────────────────────────────

func (h *AuthServiceHandler) GetJWKS(ctx context.Context, req *authnservicev1.GetJWKSRequest) (*authnservicev1.GetJWKSResponse, error) {
	resp, err := h.authService.GetJWKS(ctx, req)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return resp, nil
}
