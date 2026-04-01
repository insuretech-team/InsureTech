package models

import (
	"time"
)

// OrganisationUpdatedEvent represents a organisation_updated_event
type OrganisationUpdatedEvent struct {
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Name string `json:"name,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
