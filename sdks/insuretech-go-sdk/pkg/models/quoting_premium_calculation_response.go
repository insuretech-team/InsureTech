package models


// QuotingPremiumCalculationResponse represents a quoting_premium_calculation_response
type QuotingPremiumCalculationResponse struct {
	Calculation *PremiumCalculation `json:"calculation,omitempty"`
	Coverages []*Coverage `json:"coverages,omitempty"`
	Discounts []*Discount `json:"discounts,omitempty"`
}
