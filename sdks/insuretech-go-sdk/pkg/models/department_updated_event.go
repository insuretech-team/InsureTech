package models

import (
	"time"
)

// DepartmentUpdatedEvent represents a department_updated_event
type DepartmentUpdatedEvent struct {
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
