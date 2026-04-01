package models

import (
	"time"
)

// RoleDeletedEvent represents a role_deleted_event
type RoleDeletedEvent struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	RoleId string `json:"role_id,omitempty"`
}
