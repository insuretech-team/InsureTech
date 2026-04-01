package models

import (
	"time"
)

// AuditLog represents a audit_log
type AuditLog struct {
	Action *AuditAction `json:"action,omitempty"`
	AuditLogId string `json:"audit_log_id,omitempty"`
	Changes string `json:"changes,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	NewValues string `json:"new_values,omitempty"`
	OldValues string `json:"old_values,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TraceId string `json:"trace_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	UserId string `json:"user_id,omitempty"`
	UserRole string `json:"user_role,omitempty"`
}
