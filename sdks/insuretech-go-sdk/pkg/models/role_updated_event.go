package models

import (
	"time"
)

// RoleUpdatedEvent represents a role_updated_event
type RoleUpdatedEvent struct {
	RoleId string `json:"role_id,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
