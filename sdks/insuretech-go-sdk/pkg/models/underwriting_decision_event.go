package models

import (
	"time"
)

// UnderwritingDecisionEvent represents a underwriting_decision_event
type UnderwritingDecisionEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Decision string `json:"decision,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	RecommendedPremium *Money `json:"recommended_premium,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
