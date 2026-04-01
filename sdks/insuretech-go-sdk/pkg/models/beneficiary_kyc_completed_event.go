package models

import (
	"time"
)

// BeneficiaryKYCCompletedEvent represents a beneficiary_kyc_completed_event
type BeneficiaryKYCCompletedEvent struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	KycStatus string `json:"kyc_status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
