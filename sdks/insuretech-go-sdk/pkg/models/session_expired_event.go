package models

import (
	"time"
)

// SessionExpiredEvent represents a session_expired_event
type SessionExpiredEvent struct {
	InactivitySeconds int `json:"inactivity_seconds,omitempty"`
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}
