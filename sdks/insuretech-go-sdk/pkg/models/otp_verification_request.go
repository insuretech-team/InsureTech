package models


// OTPVerificationRequest represents a otp_verification_request
type OTPVerificationRequest struct {
	OtpId string `json:"otp_id"`
	Code string `json:"code,omitempty"`
}
