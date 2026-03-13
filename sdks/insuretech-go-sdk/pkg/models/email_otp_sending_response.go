package models


// EmailOTPSendingResponse represents a email_otp_sending_response
type EmailOTPSendingResponse struct {
	OtpId string `json:"otp_id,omitempty"`
	Message string `json:"message,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
