package models

import (
	"time"
)

// UserRole represents a user_role
type UserRole struct {
	Domain string `json:"domain"`
	AssignedBy string `json:"assigned_by,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserRoleId string `json:"user_role_id"`
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}
