package models

import (
	"time"
)

// DocumentGenerationRequestedEvent represents a document_generation_requested_event
type DocumentGenerationRequestedEvent struct {
	DocumentTemplateId string `json:"document_template_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	DocumentGenerationId string `json:"document_generation_id,omitempty"`
}
