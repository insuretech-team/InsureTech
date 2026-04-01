package models

import (
	"time"
)

// RoleRemovedEvent represents a role_removed_event
type RoleRemovedEvent struct {
	Domain string `json:"domain,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	RoleId string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
