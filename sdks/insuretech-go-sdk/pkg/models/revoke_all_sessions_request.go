package models


// RevokeAllSessionsRequest represents a revoke_all_sessions_request
type RevokeAllSessionsRequest struct {
	ExcludeCurrentSession bool `json:"exclude_current_session,omitempty"`
	Reason string `json:"reason,omitempty"`
	UserId string `json:"user_id"`
}
