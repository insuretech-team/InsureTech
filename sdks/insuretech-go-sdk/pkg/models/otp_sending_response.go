package models


// OTPSendingResponse represents a otp_sending_response
type OTPSendingResponse struct {
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
}
