package models

import (
	"time"
)

// EmployeeStatusChangedEvent represents a employee_status_changed_event
type EmployeeStatusChangedEvent struct {
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OldStatus *EmployeeStatus `json:"old_status,omitempty"`
	NewStatus *EmployeeStatus `json:"new_status,omitempty"`
	ChangedBy string `json:"changed_by,omitempty"`
}
