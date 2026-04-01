package models

import (
	"time"
)

// LeadDeletedEvent represents a lead_deleted_event
type LeadDeletedEvent struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	LeadId string `json:"lead_id,omitempty"`
	LeadName string `json:"lead_name,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
}
