package models

import (
	"time"
)

// RolePermission represents a role_permission
type RolePermission struct {
	RoleId string `json:"role_id"`
	PermissionId string `json:"permission_id"`
	AssignedAt time.Time `json:"assigned_at"`
}
