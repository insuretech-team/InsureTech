package models


// EmailUserRegistrationResponse represents a email_user_registration_response
type EmailUserRegistrationResponse struct {
	UserId string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
	VerificationEmailSent bool `json:"verification_email_sent,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	OtpExpiresInSeconds int `json:"otp_expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
