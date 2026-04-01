package models

import (
	"time"
)

// SecurityEventDetectedEvent represents a security_event_detected_event
type SecurityEventDetectedEvent struct {
	AuditEventId string `json:"audit_event_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Severity string `json:"severity,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
