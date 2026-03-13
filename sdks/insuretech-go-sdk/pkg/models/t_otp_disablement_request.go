package models


// TOTPDisablementRequest represents a t_otp_disablement_request
type TOTPDisablementRequest struct {
	UserId string `json:"user_id"`
	TotpCode string `json:"totp_code,omitempty"`
}
