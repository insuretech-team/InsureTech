package models

import (
	"time"
)

// DepartmentDeletedEvent represents a department_deleted_event
type DepartmentDeletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
