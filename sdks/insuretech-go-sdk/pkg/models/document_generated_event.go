package models

import (
	"time"
)

// DocumentGeneratedEvent represents a document_generated_event
type DocumentGeneratedEvent struct {
	EventId string `json:"event_id,omitempty"`
	DocumentGenerationId string `json:"document_generation_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
