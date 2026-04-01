package models


// LogoutResponse represents a logout_response
type LogoutResponse struct {
	SessionRevoked bool `json:"session_revoked,omitempty"`
}
