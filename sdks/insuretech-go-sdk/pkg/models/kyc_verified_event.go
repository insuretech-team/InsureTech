package models

import (
	"time"
)

// KYCVerifiedEvent represents a kyc_verified_event
type KYCVerifiedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	KycVerificationId string `json:"kyc_verification_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
