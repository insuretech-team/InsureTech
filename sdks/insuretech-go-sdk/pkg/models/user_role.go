package models

import (
	"time"
)

// UserRole represents a user_role
type UserRole struct {
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy string `json:"assigned_by,omitempty"`
	Domain string `json:"domain"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	RoleId string `json:"role_id"`
	UserId string `json:"user_id"`
	UserRoleId string `json:"user_role_id"`
}
