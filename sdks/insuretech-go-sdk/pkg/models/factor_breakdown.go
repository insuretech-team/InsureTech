package models


// FactorBreakdown represents a factor_breakdown
type FactorBreakdown struct {
	Amount float64 `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	FactorName string `json:"factor_name,omitempty"`
	FactorType string `json:"factor_type,omitempty"`
	FactorValue float64 `json:"factor_value,omitempty"`
}
