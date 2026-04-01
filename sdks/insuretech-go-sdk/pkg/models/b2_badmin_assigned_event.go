package models

import (
	"time"
)

// B2BAdminAssignedEvent represents a b2_badmin_assigned_event
type B2BAdminAssignedEvent struct {
	AssignedBy string `json:"assigned_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
