package models


// LogoutResponse represents a logout_response
type LogoutResponse struct {
	Message string `json:"message,omitempty"`
	SessionRevoked bool `json:"session_revoked,omitempty"`
	Error *Error `json:"error,omitempty"`
}
