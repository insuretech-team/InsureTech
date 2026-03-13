package models

import (
	"time"
)

// BeneficiaryCreatedEvent represents a beneficiary_created_event
type BeneficiaryCreatedEvent struct {
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	Type string `json:"type,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
}
