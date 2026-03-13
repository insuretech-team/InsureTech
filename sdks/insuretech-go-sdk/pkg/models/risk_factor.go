package models


// RiskFactor represents a risk_factor
type RiskFactor struct {
	FactorName string `json:"factor_name,omitempty"`
	ImpactScore float64 `json:"impact_score,omitempty"`
	OccurrenceCount int `json:"occurrence_count,omitempty"`
	Description string `json:"description,omitempty"`
}
