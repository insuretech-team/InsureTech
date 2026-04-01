package models


// ResetPasswordByEmailRequest represents a reset_password_by_email_request
type ResetPasswordByEmailRequest struct {
	Email string `json:"email"`
	NewPassword string `json:"new_password,omitempty"`
	OtpCode string `json:"otp_code,omitempty"`
	OtpId string `json:"otp_id"`
}
