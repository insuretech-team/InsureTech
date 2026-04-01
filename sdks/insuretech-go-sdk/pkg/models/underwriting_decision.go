package models

import (
	"time"
)

// UnderwritingDecision represents a underwriting_decision
type UnderwritingDecision struct {
	AdjustedPremium *Money `json:"adjusted_premium,omitempty"`
	AdjustmentReason string `json:"adjustment_reason,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Conditions string `json:"conditions,omitempty"`
	DecidedAt time.Time `json:"decided_at"`
	Decision *DecisionType `json:"decision"`
	Id string `json:"id"`
	Method *DecisionMethod `json:"method"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	QuoteId string `json:"quote_id"`
	Reason string `json:"reason,omitempty"`
	RiskFactors string `json:"risk_factors,omitempty"`
	RiskLevel *UnderwritingRiskLevel `json:"risk_level,omitempty"`
	RiskScore string `json:"risk_score,omitempty"`
	UnderwriterComments string `json:"underwriter_comments,omitempty"`
	UnderwriterId string `json:"underwriter_id,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
