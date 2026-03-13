package models


// ResetPasswordRequest represents a reset_password_request
type ResetPasswordRequest struct {
	MobileNumber string `json:"mobile_number"`
	OtpCode string `json:"otp_code,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}
