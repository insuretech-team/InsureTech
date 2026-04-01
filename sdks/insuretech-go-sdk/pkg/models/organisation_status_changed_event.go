package models

import (
	"time"
)

// OrganisationStatusChangedEvent represents a organisation_status_changed_event
type OrganisationStatusChangedEvent struct {
	ChangedBy string `json:"changed_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus *OrganisationStatus `json:"new_status,omitempty"`
	OldStatus *OrganisationStatus `json:"old_status,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
