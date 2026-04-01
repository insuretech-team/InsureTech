package models


// RevokeSessionRequest represents a revoke_session_request
type RevokeSessionRequest struct {
	Reason string `json:"reason,omitempty"`
	SessionId string `json:"session_id"`
}
