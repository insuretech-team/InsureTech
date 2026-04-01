package models

import (
	"time"
)

// OrganisationSuspendedEvent represents a organisation_suspended_event
type OrganisationSuspendedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	SuspendedBy string `json:"suspended_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
