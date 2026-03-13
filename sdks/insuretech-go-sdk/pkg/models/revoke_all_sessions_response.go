package models


// RevokeAllSessionsResponse represents a revoke_all_sessions_response
type RevokeAllSessionsResponse struct {
	RevokedCount int `json:"revoked_count,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
