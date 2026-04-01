package models

import (
	"time"
)

// InsurerCreatedEvent represents a insurer_created_event
type InsurerCreatedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerCode string `json:"insurer_code,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	InsurerName string `json:"insurer_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
