package models

import (
	"time"
)

// BeneficiaryRiskScoreUpdatedEvent represents a beneficiary_risk_score_updated_event
type BeneficiaryRiskScoreUpdatedEvent struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewRiskScore string `json:"new_risk_score,omitempty"`
	OldRiskScore string `json:"old_risk_score,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
