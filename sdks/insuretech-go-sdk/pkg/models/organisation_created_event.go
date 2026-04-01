package models

import (
	"time"
)

// OrganisationCreatedEvent represents a organisation_created_event
type OrganisationCreatedEvent struct {
	Code string `json:"code,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Industry string `json:"industry,omitempty"`
	Name string `json:"name,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
