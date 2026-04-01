package models


// PremiumCalculationResult represents a premium_calculation_result
type PremiumCalculationResult struct {
	BasePremium float64 `json:"base_premium,omitempty"`
	Currency string `json:"currency,omitempty"`
	FactorBreakdown []*FactorBreakdown `json:"factor_breakdown,omitempty"`
	GrossPremium float64 `json:"gross_premium,omitempty"`
	NetPremium float64 `json:"net_premium,omitempty"`
	TotalDiscounts float64 `json:"total_discounts,omitempty"`
	TotalLoadings float64 `json:"total_loadings,omitempty"`
}
