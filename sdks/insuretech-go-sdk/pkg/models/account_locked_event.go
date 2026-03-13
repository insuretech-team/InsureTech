package models

import (
	"time"
)

// AccountLockedEvent represents a account_locked_event
type AccountLockedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	LockedUntil time.Time `json:"locked_until,omitempty"`
}
