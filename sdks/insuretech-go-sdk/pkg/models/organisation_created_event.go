package models

import (
	"time"
)

// OrganisationCreatedEvent represents a organisation_created_event
type OrganisationCreatedEvent struct {
	Industry string `json:"industry,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}
