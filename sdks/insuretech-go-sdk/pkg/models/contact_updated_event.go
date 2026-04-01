package models

import (
	"time"
)

// ContactUpdatedEvent represents a contact_updated_event
type ContactUpdatedEvent struct {
	ChangedFields []string `json:"changed_fields,omitempty"`
	ContactId string `json:"contact_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
