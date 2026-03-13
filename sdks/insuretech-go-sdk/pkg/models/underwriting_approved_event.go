package models

import (
	"time"
)

// UnderwritingApprovedEvent represents a underwriting_approved_event
type UnderwritingApprovedEvent struct {
	DecisionMethod string `json:"decision_method,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	DecisionId string `json:"decision_id,omitempty"`
}
