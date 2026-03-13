package models

import (
	"time"
)

// User represents a user
type User struct {
	UserType *UserType `json:"user_type,omitempty"`
	TotpEnabled bool `json:"totp_enabled,omitempty"`
	EmailVerified bool `json:"email_verified,omitempty"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
	BiometricTokenEnc string `json:"biometric_token_enc,omitempty"`
	ActivePoliciesCount int `json:"active_policies_count,omitempty"`
	PendingClaimsCount int `json:"pending_claims_count,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PreferredAuthMethod string `json:"preferred_auth_method,omitempty"`
	WalletPaymentMethod string `json:"wallet_payment_method,omitempty"`
	MobileNumberIdx string `json:"mobile_number_idx,omitempty"`
	Status *UserStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	LastLoginSessionType string `json:"last_login_session_type,omitempty"`
	EmailLockedUntil time.Time `json:"email_locked_until,omitempty"`
	BiometricTokenIdx string `json:"biometric_token_idx,omitempty"`
	EmailIdx string `json:"email_idx,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	EmailLoginAttempts int `json:"email_login_attempts,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email,omitempty"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	Username string `json:"username,omitempty"`
	WalletBalance *Money `json:"wallet_balance,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	TotpSecretEnc string `json:"totp_secret_enc,omitempty"`
	LoginAttempts int `json:"login_attempts,omitempty"`
	NotificationPreference string `json:"notification_preference,omitempty"`
	UserId string `json:"user_id,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	LockedUntil time.Time `json:"locked_until,omitempty"`
}
