package models


// OTPVerificationRequest represents a otp_verification_request
type OTPVerificationRequest struct {
	Code string `json:"code,omitempty"`
	DeviceId string `json:"device_id"`
	DeviceType string `json:"device_type,omitempty"`
	OtpId string `json:"otp_id"`
}
