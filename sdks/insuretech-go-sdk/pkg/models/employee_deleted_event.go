package models

import (
	"time"
)

// EmployeeDeletedEvent represents a employee_deleted_event
type EmployeeDeletedEvent struct {
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
