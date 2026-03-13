package models

import (
	"time"
)

// DocumentGenerationFailedEvent represents a document_generation_failed_event
type DocumentGenerationFailedEvent struct {
	EventId string `json:"event_id,omitempty"`
	DocumentGenerationId string `json:"document_generation_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
