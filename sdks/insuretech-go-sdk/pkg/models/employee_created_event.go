package models

import (
	"time"
)

// EmployeeCreatedEvent represents a employee_created_event
type EmployeeCreatedEvent struct {
	CreatedBy string `json:"created_by,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EmployeeId string `json:"employee_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Name string `json:"name,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
