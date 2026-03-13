package models

import (
	"time"
)

// AuditEvent represents a audit_event
type AuditEvent struct {
	Category *EventCategory `json:"category"`
	EventType string `json:"event_type"`
	Severity *EventSeverity `json:"severity"`
	Description string `json:"description"`
	EntityId string `json:"entity_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Id string `json:"id"`
	UserId string `json:"user_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}
