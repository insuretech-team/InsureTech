package models

import (
	"time"
)

// BeneficiaryRiskScoreUpdatedEvent represents a beneficiary_risk_score_updated_event
type BeneficiaryRiskScoreUpdatedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	OldRiskScore string `json:"old_risk_score,omitempty"`
	NewRiskScore string `json:"new_risk_score,omitempty"`
	Reason string `json:"reason,omitempty"`
}
