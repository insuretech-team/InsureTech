package models

import (
	"time"
)

// RoleRevokedEvent represents a role_revoked_event
type RoleRevokedEvent struct {
	AssignmentId string `json:"assignment_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	RoleId string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	RevokedBy string `json:"revoked_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
