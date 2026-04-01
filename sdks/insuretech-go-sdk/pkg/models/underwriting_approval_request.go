package models


// UnderwritingApprovalRequest represents a underwriting_approval_request
type UnderwritingApprovalRequest struct {
	AdjustedPremium *Money `json:"adjusted_premium,omitempty"`
	Comments string `json:"comments,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	PremiumAdjusted bool `json:"premium_adjusted,omitempty"`
	QuoteId string `json:"quote_id"`
	RiskLevel string `json:"risk_level,omitempty"`
	UnderwriterId string `json:"underwriter_id"`
}
