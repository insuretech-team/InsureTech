package models

import (
	"time"
)

// User represents a user
type User struct {
	ActivePoliciesCount int `json:"active_policies_count,omitempty"`
	BiometricTokenEnc string `json:"biometric_token_enc,omitempty"`
	BiometricTokenIdx string `json:"biometric_token_idx,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Email string `json:"email,omitempty"`
	EmailIdx string `json:"email_idx,omitempty"`
	EmailLockedUntil time.Time `json:"email_locked_until,omitempty"`
	EmailLoginAttempts int `json:"email_login_attempts,omitempty"`
	EmailVerified bool `json:"email_verified,omitempty"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
	LastLoginSessionType string `json:"last_login_session_type,omitempty"`
	LockedUntil time.Time `json:"locked_until,omitempty"`
	LoginAttempts int `json:"login_attempts,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	MobileNumberIdx string `json:"mobile_number_idx,omitempty"`
	NotificationPreference string `json:"notification_preference,omitempty"`
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	PendingClaimsCount int `json:"pending_claims_count,omitempty"`
	PreferredAuthMethod string `json:"preferred_auth_method,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	Status *UserStatus `json:"status,omitempty"`
	TotpEnabled bool `json:"totp_enabled,omitempty"`
	TotpSecretEnc string `json:"totp_secret_enc,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	UserId string `json:"user_id,omitempty"`
	UserType *UserType `json:"user_type,omitempty"`
	Username string `json:"username,omitempty"`
	WalletBalance *Money `json:"wallet_balance,omitempty"`
	WalletPaymentMethod string `json:"wallet_payment_method,omitempty"`
}
