package models

import (
	"time"
)

// SecurityEventDetectedEvent represents a security_event_detected_event
type SecurityEventDetectedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AuditEventId string `json:"audit_event_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	Severity string `json:"severity,omitempty"`
	UserId string `json:"user_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
