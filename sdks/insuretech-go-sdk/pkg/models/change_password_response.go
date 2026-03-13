package models


// ChangePasswordResponse represents a change_password_response
type ChangePasswordResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
