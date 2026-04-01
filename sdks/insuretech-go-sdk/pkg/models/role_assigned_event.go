package models

import (
	"time"
)

// RoleAssignedEvent represents a role_assigned_event
type RoleAssignedEvent struct {
	AssignedBy string `json:"assigned_by,omitempty"`
	Domain string `json:"domain,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	RoleId string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
