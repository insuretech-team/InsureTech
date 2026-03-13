package models


// ResendOTPRequest represents a resend_otp_request
type ResendOTPRequest struct {
	OriginalOtpId string `json:"original_otp_id"`
	Reason string `json:"reason,omitempty"`
}
