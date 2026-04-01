package models

import (
	"time"
)

// LeadUpdatedEvent represents a lead_updated_event
type LeadUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	LeadId string `json:"lead_id,omitempty"`
	NewStatus *LeadStatus `json:"new_status,omitempty"`
	OldStatus *LeadStatus `json:"old_status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
