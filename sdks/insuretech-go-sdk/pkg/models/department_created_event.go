package models

import (
	"time"
)

// DepartmentCreatedEvent represents a department_created_event
type DepartmentCreatedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Name string `json:"name,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}
