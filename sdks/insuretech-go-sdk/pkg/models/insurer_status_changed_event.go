package models

import (
	"time"
)

// InsurerStatusChangedEvent represents a insurer_status_changed_event
type InsurerStatusChangedEvent struct {
	InsurerId string `json:"insurer_id,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	Reason string `json:"reason,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
