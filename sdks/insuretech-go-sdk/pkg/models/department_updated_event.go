package models

import (
	"time"
)

// DepartmentUpdatedEvent represents a department_updated_event
type DepartmentUpdatedEvent struct {
	DepartmentId string `json:"department_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
