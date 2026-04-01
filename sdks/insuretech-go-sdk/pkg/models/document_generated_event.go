package models

import (
	"time"
)

// DocumentGeneratedEvent represents a document_generated_event
type DocumentGeneratedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	DocumentGenerationId string `json:"document_generation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
