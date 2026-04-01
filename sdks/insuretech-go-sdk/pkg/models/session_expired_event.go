package models

import (
	"time"
)

// SessionExpiredEvent represents a session_expired_event
type SessionExpiredEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	InactivitySeconds int `json:"inactivity_seconds,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
