package models

import (
	"time"
)

// SessionRevokedEvent represents a session_revoked_event
type SessionRevokedEvent struct {
	EventId string `json:"event_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	RevokedBy string `json:"revoked_by,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
