package models

import (
	"time"
)

// RoleAssignmentRequest represents a role_assignment_request
type RoleAssignmentRequest struct {
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
	Domain string `json:"domain,omitempty"`
	AssignedBy string `json:"assigned_by,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
