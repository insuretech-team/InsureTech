package models


// ResetPasswordByEmailResponse represents a reset_password_by_email_response
type ResetPasswordByEmailResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
