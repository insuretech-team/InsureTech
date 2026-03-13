package models

import (
	"time"
)

// PermissionCreatedEvent represents a permission_created_event
type PermissionCreatedEvent struct {
	PermissionName string `json:"permission_name,omitempty"`
	Resource string `json:"resource,omitempty"`
	Action string `json:"action,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PermissionId string `json:"permission_id,omitempty"`
}
