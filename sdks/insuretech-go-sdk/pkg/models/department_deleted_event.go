package models

import (
	"time"
)

// DepartmentDeletedEvent represents a department_deleted_event
type DepartmentDeletedEvent struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
