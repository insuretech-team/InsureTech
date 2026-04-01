package models


// EmailUserRegistrationResponse represents a email_user_registration_response
type EmailUserRegistrationResponse struct {
	OtpExpiresInSeconds int `json:"otp_expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	VerificationEmailSent bool `json:"verification_email_sent,omitempty"`
}
