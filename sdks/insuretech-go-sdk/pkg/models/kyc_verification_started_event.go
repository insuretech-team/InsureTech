package models

import (
	"time"
)

// KYCVerificationStartedEvent represents a kyc_verification_started_event
type KYCVerificationStartedEvent struct {
	EntityId string `json:"entity_id,omitempty"`
	Method string `json:"method,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	KycVerificationId string `json:"kyc_verification_id,omitempty"`
	Type string `json:"type,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
}
