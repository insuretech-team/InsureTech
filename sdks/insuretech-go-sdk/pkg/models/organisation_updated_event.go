package models

import (
	"time"
)

// OrganisationUpdatedEvent represents a organisation_updated_event
type OrganisationUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Name string `json:"name,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
