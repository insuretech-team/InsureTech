package models

import (
	"time"
)

// OrganisationApprovedEvent represents a organisation_approved_event
type OrganisationApprovedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
