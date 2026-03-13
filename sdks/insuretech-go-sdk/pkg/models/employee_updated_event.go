package models

import (
	"time"
)

// EmployeeUpdatedEvent represents a employee_updated_event
type EmployeeUpdatedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
}
