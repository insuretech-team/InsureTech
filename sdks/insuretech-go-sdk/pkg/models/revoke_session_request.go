package models


// RevokeSessionRequest represents a revoke_session_request
type RevokeSessionRequest struct {
	SessionId string `json:"session_id"`
	Reason string `json:"reason,omitempty"`
}
