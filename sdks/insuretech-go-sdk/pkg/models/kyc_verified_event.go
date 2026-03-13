package models

import (
	"time"
)

// KYCVerifiedEvent represents a kyc_verified_event
type KYCVerifiedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	KycVerificationId string `json:"kyc_verification_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
