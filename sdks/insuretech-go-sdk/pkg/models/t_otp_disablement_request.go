package models


// TOTPDisablementRequest represents a t_otp_disablement_request
type TOTPDisablementRequest struct {
	TotpCode string `json:"totp_code,omitempty"`
	UserId string `json:"user_id"`
}
