package models

import (
	"time"
)

// InsurerConfigUpdatedEvent represents a insurer_config_updated_event
type InsurerConfigUpdatedEvent struct {
	InsurerId string `json:"insurer_id,omitempty"`
	ConfigType string `json:"config_type,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
