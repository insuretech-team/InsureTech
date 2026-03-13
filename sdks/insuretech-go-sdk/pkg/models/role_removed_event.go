package models

import (
	"time"
)

// RoleRemovedEvent represents a role_removed_event
type RoleRemovedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	RoleId string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Domain string `json:"domain,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
