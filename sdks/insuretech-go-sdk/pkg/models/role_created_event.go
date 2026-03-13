package models

import (
	"time"
)

// RoleCreatedEvent represents a role_created_event
type RoleCreatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	RoleId string `json:"role_id,omitempty"`
	Name string `json:"name,omitempty"`
	Portal *Portal `json:"portal,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
