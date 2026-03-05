package domain

import (
	"context"
	"time"

	entity "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
)

// Primary Port (Inbound) - Mobile/SMS auth (B2C Customer + Agent)
type AuthService interface {
	// Mobile OTP auth (B2C + Agent)
	Login(ctx context.Context, req *authnservicev1.LoginRequest) (*authnservicev1.LoginResponse, error)
	Register(ctx context.Context, req *authnservicev1.RegisterRequest) (*authnservicev1.RegisterResponse, error)
	SendOTP(ctx context.Context, req *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error)
	VerifyOTP(ctx context.Context, req *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error)
	RefreshToken(ctx context.Context, req *authnservicev1.RefreshTokenRequest) (*authnservicev1.RefreshTokenResponse, error)
	Logout(ctx context.Context, req *authnservicev1.LogoutRequest) (*authnservicev1.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *authnservicev1.ChangePasswordRequest) (*authnservicev1.ChangePasswordResponse, error)
	ResetPassword(ctx context.Context, req *authnservicev1.ResetPasswordRequest) (*authnservicev1.ResetPasswordResponse, error)
	ValidateToken(ctx context.Context, req *authnservicev1.ValidateTokenRequest) (*authnservicev1.ValidateTokenResponse, error)
	RevokeAllSessions(ctx context.Context, req *authnservicev1.RevokeAllSessionsRequest) (*authnservicev1.RevokeAllSessionsResponse, error)

	// Email OTP auth (Business Beneficiary + System User — web portal only)
	RegisterEmailUser(ctx context.Context, req *authnservicev1.RegisterEmailUserRequest) (*authnservicev1.RegisterEmailUserResponse, error)
	SendEmailOTP(ctx context.Context, req *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error)
	VerifyEmail(ctx context.Context, req *authnservicev1.VerifyEmailRequest) (*authnservicev1.VerifyEmailResponse, error)
	EmailLogin(ctx context.Context, req *authnservicev1.EmailLoginRequest) (*authnservicev1.EmailLoginResponse, error)
	RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error)
	ResetPasswordByEmail(ctx context.Context, req *authnservicev1.ResetPasswordByEmailRequest) (*authnservicev1.ResetPasswordByEmailResponse, error)
}

// Secondary Port (Outbound - DB) - Session storage
type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByID(ctx context.Context, id string) (*entity.Session, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.Session, error)
	UpdateLastActivity(ctx context.Context, sessionID string) error
	UpdateTokens(ctx context.Context, sessionID, newAccessJTI, newRefreshJTI string) error
	Revoke(ctx context.Context, id string) error
	RevokeAllByUserID(ctx context.Context, userID string, excludeSessionID string) error
	ListByUserID(ctx context.Context, userID string, activeOnly bool, sessionType *entity.SessionType) ([]*entity.Session, error)
	CleanupExpiredSessions(ctx context.Context) (int64, error)
}

// Secondary Port (Outbound - DB) - User storage
type UserRepository interface {
	// Core CRUD
	Create(ctx context.Context, mobile, passwordHash, email string, status entity.UserStatus) (*entity.User, error)
	CreateFull(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByMobileNumber(ctx context.Context, mobile string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Updates
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateStatus(ctx context.Context, userID string, status entity.UserStatus) error
	UpdateLastLogin(ctx context.Context, userID, sessionType string) error

	// Email auth specific
	UpdateEmailVerified(ctx context.Context, userID string) error
	IncrementEmailLoginAttempts(ctx context.Context, userID string) (int32, error)
	LockEmailAuth(ctx context.Context, userID string, lockDuration time.Duration) error
	ResetEmailLoginAttempts(ctx context.Context, userID string) error
}

// Secondary Port (Outbound - Messaging) - Event publishing
type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, userID, mobile, email, ip, deviceType string) error
	PublishUserLoggedIn(ctx context.Context, userID, sessionID, sessionType, ip, deviceType, userAgent string) error
	PublishUserLoggedOut(ctx context.Context, userID, sessionID, sessionType, reason, ip, deviceType string) error
	PublishLoginFailed(ctx context.Context, userID, mobile, reason, ip, deviceType, userAgent string, failedAttempts int32) error
	PublishTokenRefreshed(ctx context.Context, userID, sessionID, oldJTI, newAccessJTI, newRefreshJTI, ip, deviceType, userAgent string) error
	PublishPasswordChanged(ctx context.Context, userID, ip, changedBy string) error
	PublishPasswordResetRequested(ctx context.Context, userID, mobile, ip, deviceType string) error
	PublishSessionRevoked(ctx context.Context, userID, sessionID, sessionType, revokedBy, reason string) error
	PublishCSRFValidationFailed(ctx context.Context, userID, sessionID, expectedHash, receivedHash, ip, userAgent, path, method string) error
	PublishOTPSent(ctx context.Context, otpID, recipientMasked, otpType, channel, provider, senderID, providerMessageID string, maskingUsed bool) error
	PublishOTPVerified(ctx context.Context, otpID, userID string, attempts int32) error
	PublishSMSDeliveryReport(ctx context.Context, otpID, providerMsgID, msisdnMasked, status, errorCode, carrier string, deliveredAt time.Time) error
	PublishEmailVerificationSent(ctx context.Context, userID, email, otpID, otpType, ipAddress string) error
	PublishEmailVerified(ctx context.Context, userID, email string) error
	PublishEmailLoginSucceeded(ctx context.Context, userID, sessionID, email, userType, ipAddress, userAgent, deviceName string) error
	PublishEmailLoginFailed(ctx context.Context, email, reason string, attempts int32, ipAddress, userAgent string) error
	PublishPasswordResetByEmailRequested(ctx context.Context, userID, email, otpID, ipAddress string) error
}
