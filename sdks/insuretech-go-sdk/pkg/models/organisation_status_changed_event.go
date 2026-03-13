package models

import (
	"time"
)

// OrganisationStatusChangedEvent represents a organisation_status_changed_event
type OrganisationStatusChangedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OldStatus *OrganisationStatus `json:"old_status,omitempty"`
	NewStatus *OrganisationStatus `json:"new_status,omitempty"`
	ChangedBy string `json:"changed_by,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
