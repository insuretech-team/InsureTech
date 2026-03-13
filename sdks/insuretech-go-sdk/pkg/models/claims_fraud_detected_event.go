package models

import (
	"time"
)

// ClaimsFraudDetectedEvent represents a claims_fraud_detected_event
type ClaimsFraudDetectedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	RiskIndicators []string `json:"risk_indicators,omitempty"`
}
