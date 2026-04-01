package models

import (
	"time"
)

// RoleAssignmentRequest represents a role_assignment_request
type RoleAssignmentRequest struct {
	AssignedBy string `json:"assigned_by,omitempty"`
	Domain string `json:"domain,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	RoleId string `json:"role_id"`
	UserId string `json:"user_id"`
}
