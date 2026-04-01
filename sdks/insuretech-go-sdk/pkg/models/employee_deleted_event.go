package models

import (
	"time"
)

// EmployeeDeletedEvent represents a employee_deleted_event
type EmployeeDeletedEvent struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
