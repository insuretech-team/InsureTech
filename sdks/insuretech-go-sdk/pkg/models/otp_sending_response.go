package models


// OTPSendingResponse represents a otp_sending_response
type OTPSendingResponse struct {
	OtpId string `json:"otp_id,omitempty"`
	Message string `json:"message,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
