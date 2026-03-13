package models

import (
	"time"
)

// AuditLogCreatedEvent represents a audit_log_created_event
type AuditLogCreatedEvent struct {
	UserId string `json:"user_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AuditLogId string `json:"audit_log_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	Action string `json:"action,omitempty"`
}
