package models

import (
	"time"
)

// KYCRejectedEvent represents a kyc_rejected_event
type KYCRejectedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	KycVerificationId string `json:"kyc_verification_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	Reason string `json:"reason,omitempty"`
}
