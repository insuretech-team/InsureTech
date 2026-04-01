package models


// QuotingPremiumCalculationRequest represents a quoting_premium_calculation_request
type QuotingPremiumCalculationRequest struct {
	Parameters *QuoteParameters `json:"parameters,omitempty"`
	ProductId string `json:"product_id"`
}
