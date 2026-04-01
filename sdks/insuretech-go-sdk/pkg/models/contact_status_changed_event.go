package models

import (
	"time"
)

// ContactStatusChangedEvent represents a contact_status_changed_event
type ContactStatusChangedEvent struct {
	ChangedAt time.Time `json:"changed_at,omitempty"`
	ContactId string `json:"contact_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus *ContactStatus `json:"new_status,omitempty"`
	OldStatus *ContactStatus `json:"old_status,omitempty"`
	Reason string `json:"reason,omitempty"`
}
