package models


// PremiumCalculationInput represents a premium_calculation_input
type PremiumCalculationInput struct {
	CoveragePeriodMonths int `json:"coverage_period_months,omitempty"`
	CoverageType string `json:"coverage_type,omitempty"`
	Discounts []string `json:"discounts,omitempty"`
	Loadings []string `json:"loadings,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	RatingFactors map[string]interface{} `json:"rating_factors,omitempty"`
	RiskCharacteristics map[string]interface{} `json:"risk_characteristics,omitempty"`
	SumInsured float64 `json:"sum_insured,omitempty"`
}
