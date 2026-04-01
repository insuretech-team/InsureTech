package models

import (
	"time"
)

// InsurerConfigUpdatedEvent represents a insurer_config_updated_event
type InsurerConfigUpdatedEvent struct {
	ConfigType string `json:"config_type,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
