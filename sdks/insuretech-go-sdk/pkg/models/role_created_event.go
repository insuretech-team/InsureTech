package models

import (
	"time"
)

// RoleCreatedEvent represents a role_created_event
type RoleCreatedEvent struct {
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Name string `json:"name,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Portal *Portal `json:"portal,omitempty"`
	RoleId string `json:"role_id,omitempty"`
}
