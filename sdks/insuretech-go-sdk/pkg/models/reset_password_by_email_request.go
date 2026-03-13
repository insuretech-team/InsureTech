package models


// ResetPasswordByEmailRequest represents a reset_password_by_email_request
type ResetPasswordByEmailRequest struct {
	OtpId string `json:"otp_id"`
	OtpCode string `json:"otp_code,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
	Email string `json:"email"`
}
