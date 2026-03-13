package models

import (
	"time"
)

// EmployeeCreatedEvent represents a employee_created_event
type EmployeeCreatedEvent struct {
	DepartmentId string `json:"department_id,omitempty"`
	Name string `json:"name,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	EmployeeId string `json:"employee_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
