package models


// UnderwritingApprovalRequest represents a underwriting_approval_request
type UnderwritingApprovalRequest struct {
	UnderwriterId string `json:"underwriter_id"`
	RiskLevel string `json:"risk_level,omitempty"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	AdjustedPremium *Money `json:"adjusted_premium,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Comments string `json:"comments,omitempty"`
	QuoteId string `json:"quote_id"`
}
