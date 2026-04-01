package models

import (
	"time"
)

// DocumentGenerationRequestedEvent represents a document_generation_requested_event
type DocumentGenerationRequestedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	DocumentGenerationId string `json:"document_generation_id,omitempty"`
	DocumentTemplateId string `json:"document_template_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
