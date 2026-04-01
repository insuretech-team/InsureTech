package models


// EmailOTPSendingResponse represents a email_otp_sending_response
type EmailOTPSendingResponse struct {
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
}
