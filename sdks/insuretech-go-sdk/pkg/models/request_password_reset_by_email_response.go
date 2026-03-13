package models


// RequestPasswordResetByEmailResponse represents a request_password_reset_by_email_response
type RequestPasswordResetByEmailResponse struct {
	Message string `json:"message,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
}
