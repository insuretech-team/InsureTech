package models

import (
	"time"
)

// HealthDeclarationSubmittedEvent represents a health_declaration_submitted_event
type HealthDeclarationSubmittedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	HasPreExistingConditions bool `json:"has_pre_existing_conditions,omitempty"`
	MedicalExamRequired bool `json:"medical_exam_required,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
