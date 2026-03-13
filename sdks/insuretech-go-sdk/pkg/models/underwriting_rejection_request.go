package models


// UnderwritingRejectionRequest represents a underwriting_rejection_request
type UnderwritingRejectionRequest struct {
	UnderwriterId string `json:"underwriter_id"`
	Reason string `json:"reason,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	Comments string `json:"comments,omitempty"`
	QuoteId string `json:"quote_id"`
}
