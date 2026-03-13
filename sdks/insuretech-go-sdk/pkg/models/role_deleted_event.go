package models

import (
	"time"
)

// RoleDeletedEvent represents a role_deleted_event
type RoleDeletedEvent struct {
	RoleId string `json:"role_id,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
