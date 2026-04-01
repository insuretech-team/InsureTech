package models

import (
	"time"
)

// ContactDeletedEvent represents a contact_deleted_event
type ContactDeletedEvent struct {
	ContactId string `json:"contact_id,omitempty"`
	ContactName string `json:"contact_name,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
}
