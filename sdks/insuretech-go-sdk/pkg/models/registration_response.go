package models


// RegistrationResponse represents a registration_response
type RegistrationResponse struct {
	OtpExpiresInSeconds int `json:"otp_expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	OtpSent bool `json:"otp_sent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
