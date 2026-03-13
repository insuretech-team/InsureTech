package models


// RevokeSessionResponse represents a revoke_session_response
type RevokeSessionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
