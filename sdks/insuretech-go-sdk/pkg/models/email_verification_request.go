package models


// EmailVerificationRequest represents a email_verification_request
type EmailVerificationRequest struct {
	Code string `json:"code,omitempty"`
	OtpId string `json:"otp_id"`
}
