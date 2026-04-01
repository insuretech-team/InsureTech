package models


// EmailLoginRequest represents a email_login_request
type EmailLoginRequest struct {
	Code string `json:"code,omitempty"`
	DeviceId string `json:"device_id"`
	DeviceName string `json:"device_name,omitempty"`
	Email string `json:"email"`
	OtpId string `json:"otp_id"`
}
