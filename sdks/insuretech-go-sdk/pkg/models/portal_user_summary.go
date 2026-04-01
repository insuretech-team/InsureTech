package models


// PortalUserSummary represents a portal_user_summary
type PortalUserSummary struct {
	Email string `json:"email,omitempty"`
	EmailVerified bool `json:"email_verified,omitempty"`
	FullName string `json:"full_name,omitempty"`
	KycVerified bool `json:"kyc_verified,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	UserId string `json:"user_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
}
