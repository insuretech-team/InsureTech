package models

import (
	"time"
)

// BeneficiaryCreatedEvent represents a beneficiary_created_event
type BeneficiaryCreatedEvent struct {
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
