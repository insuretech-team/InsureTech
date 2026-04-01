package models

import (
	"time"
)

// AuditInfo represents a audit_info
type AuditInfo struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
