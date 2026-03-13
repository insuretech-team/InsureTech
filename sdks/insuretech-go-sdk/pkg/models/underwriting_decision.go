package models

import (
	"time"
)

// UnderwritingDecision represents a underwriting_decision
type UnderwritingDecision struct {
	Id string `json:"id"`
	QuoteId string `json:"quote_id"`
	UnderwriterComments string `json:"underwriter_comments,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Decision *DecisionType `json:"decision"`
	Conditions string `json:"conditions,omitempty"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	AdjustedPremium *Money `json:"adjusted_premium,omitempty"`
	Method *DecisionMethod `json:"method"`
	RiskLevel *UnderwritingRiskLevel `json:"risk_level,omitempty"`
	Reason string `json:"reason,omitempty"`
	UnderwriterId string `json:"underwriter_id,omitempty"`
	DecidedAt time.Time `json:"decided_at"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
	RiskScore string `json:"risk_score,omitempty"`
	RiskFactors string `json:"risk_factors,omitempty"`
	AdjustmentReason string `json:"adjustment_reason,omitempty"`
}
