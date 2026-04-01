package models

import (
	"time"
)

// Role represents a role
type Role struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active"`
	IsSystem bool `json:"is_system"`
	Name string `json:"name"`
	Portal *Portal `json:"portal"`
	RoleId string `json:"role_id"`
	UpdatedAt time.Time `json:"updated_at"`
}
