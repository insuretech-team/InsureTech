package models

import (
	"time"
)

// FraudConfirmedEvent represents a fraud_confirmed_event
type FraudConfirmedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FraudCaseId string `json:"fraud_case_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
