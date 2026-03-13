package models

import (
	"time"
)

// B2BAdminAssignedEvent represents a b2_badmin_assigned_event
type B2BAdminAssignedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	AssignedBy string `json:"assigned_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
