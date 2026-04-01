package models

import (
	"time"
)

// AuditEvent represents a audit_event
type AuditEvent struct {
	Category *EventCategory `json:"category"`
	Description string `json:"description"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventType string `json:"event_type"`
	Id string `json:"id"`
	Metadata string `json:"metadata,omitempty"`
	Severity *EventSeverity `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
	UserId string `json:"user_id,omitempty"`
}
