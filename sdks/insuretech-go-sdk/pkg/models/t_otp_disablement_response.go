package models


// TOTPDisablementResponse represents a t_otp_disablement_response
type TOTPDisablementResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
