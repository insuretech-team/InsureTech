package models

import (
	"time"
)

// Permission represents a permission
type Permission struct {
	PermissionId string `json:"permission_id"`
	PermissionName string `json:"permission_name"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
