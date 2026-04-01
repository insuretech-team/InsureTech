package models

import (
	"time"
)

// PasswordChangedEvent represents a password_changed_event
type PasswordChangedEvent struct {
	ChangedBy string `json:"changed_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
