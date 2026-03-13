package models

import (
	"time"
)

// OrganisationSuspendedEvent represents a organisation_suspended_event
type OrganisationSuspendedEvent struct {
	SuspendedBy string `json:"suspended_by,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
