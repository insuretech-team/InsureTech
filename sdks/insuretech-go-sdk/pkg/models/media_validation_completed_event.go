package models

import (
	"time"
)

// MediaValidationCompletedEvent represents a media_validation_completed_event
type MediaValidationCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	ValidationErrors []string `json:"validation_errors,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
}
