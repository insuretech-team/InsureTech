package models


// RegistrationResponse represents a registration_response
type RegistrationResponse struct {
	OtpId string `json:"otp_id,omitempty"`
	OtpExpiresInSeconds int `json:"otp_expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
	OtpSent bool `json:"otp_sent,omitempty"`
}
