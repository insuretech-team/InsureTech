package models


// RequestPasswordResetByEmailResponse represents a request_password_reset_by_email_response
type RequestPasswordResetByEmailResponse struct {
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
}
