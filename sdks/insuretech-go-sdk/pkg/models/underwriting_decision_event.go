package models

import (
	"time"
)

// UnderwritingDecisionEvent represents a underwriting_decision_event
type UnderwritingDecisionEvent struct {
	RecommendedPremium *Money `json:"recommended_premium,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Decision string `json:"decision,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
}
