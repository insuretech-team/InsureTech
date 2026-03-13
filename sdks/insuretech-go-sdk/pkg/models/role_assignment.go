package models

import (
	"time"
)

// RoleAssignment represents a role_assignment
type RoleAssignment struct {
	AssignedBy string `json:"assigned_by"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	AssignmentId string `json:"assignment_id"`
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
	TenantId string `json:"tenant_id,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
}
