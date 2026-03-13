package models

import (
	"time"
)

// PermissionCheckedEvent represents a permission_checked_event
type PermissionCheckedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Resource string `json:"resource,omitempty"`
	Action string `json:"action,omitempty"`
	Granted bool `json:"granted,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
