package models

import (
	"time"
)

// InsurerStatusChangedEvent represents a insurer_status_changed_event
type InsurerStatusChangedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
