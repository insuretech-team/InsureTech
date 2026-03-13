package models


// EmailVerificationResponse represents a email_verification_response
type EmailVerificationResponse struct {
	Verified bool `json:"verified,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
