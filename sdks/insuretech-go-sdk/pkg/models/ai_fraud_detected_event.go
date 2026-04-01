package models

import (
	"time"
)

// AiFraudDetectedEvent represents a ai_fraud_detected_event
type AiFraudDetectedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	RequiresManualReview bool `json:"requires_manual_review,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	Severity string `json:"severity,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
