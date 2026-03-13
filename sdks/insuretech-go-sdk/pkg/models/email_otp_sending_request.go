package models


// EmailOTPSendingRequest represents a email_otp_sending_request
type EmailOTPSendingRequest struct {
	Email string `json:"email"`
	Type string `json:"type"`
}
