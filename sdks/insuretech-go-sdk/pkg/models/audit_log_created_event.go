package models

import (
	"time"
)

// AuditLogCreatedEvent represents a audit_log_created_event
type AuditLogCreatedEvent struct {
	Action string `json:"action,omitempty"`
	AuditLogId string `json:"audit_log_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
