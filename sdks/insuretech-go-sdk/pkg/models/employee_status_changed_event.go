package models

import (
	"time"
)

// EmployeeStatusChangedEvent represents a employee_status_changed_event
type EmployeeStatusChangedEvent struct {
	ChangedBy string `json:"changed_by,omitempty"`
	EmployeeUuid string `json:"employee_uuid,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus *EmployeeStatus `json:"new_status,omitempty"`
	OldStatus *EmployeeStatus `json:"old_status,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
