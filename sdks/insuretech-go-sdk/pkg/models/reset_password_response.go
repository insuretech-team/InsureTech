package models


// ResetPasswordResponse represents a reset_password_response
type ResetPasswordResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
