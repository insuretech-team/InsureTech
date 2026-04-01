package models

import (
	"time"
)

// DepartmentCreatedEvent represents a department_created_event
type DepartmentCreatedEvent struct {
	CreatedBy string `json:"created_by,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Name string `json:"name,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
