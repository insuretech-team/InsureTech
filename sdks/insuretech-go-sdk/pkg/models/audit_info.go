package models

import (
	"time"
)

// AuditInfo represents a audit_info
type AuditInfo struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
