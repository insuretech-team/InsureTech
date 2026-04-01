package models


// EmailVerificationResponse represents a email_verification_response
type EmailVerificationResponse struct {
	UserId string `json:"user_id,omitempty"`
	Verified bool `json:"verified,omitempty"`
}
