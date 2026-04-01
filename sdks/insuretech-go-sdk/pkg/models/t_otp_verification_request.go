package models


// TOTPVerificationRequest represents a t_otp_verification_request
type TOTPVerificationRequest struct {
	MfaSessionToken string `json:"mfa_session_token,omitempty"`
	TotpCode string `json:"totp_code,omitempty"`
	UserId string `json:"user_id"`
}
