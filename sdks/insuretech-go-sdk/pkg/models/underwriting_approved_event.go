package models

import (
	"time"
)

// UnderwritingApprovedEvent represents a underwriting_approved_event
type UnderwritingApprovedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	DecisionId string `json:"decision_id,omitempty"`
	DecisionMethod string `json:"decision_method,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
