package models

import (
	"time"
)

// Role represents a role
type Role struct {
	Name string `json:"name"`
	IsSystem bool `json:"is_system"`
	IsActive bool `json:"is_active"`
	CreatedBy string `json:"created_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Portal *Portal `json:"portal"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	RoleId string `json:"role_id"`
}
