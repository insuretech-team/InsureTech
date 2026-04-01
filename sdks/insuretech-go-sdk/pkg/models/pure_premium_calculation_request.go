package models


// PurePremiumCalculationRequest represents a pure_premium_calculation_request
type PurePremiumCalculationRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	ClaimSeverity float64 `json:"claim_severity,omitempty"`
	ExpectedClaims float64 `json:"expected_claims,omitempty"`
	ExposureUnits float64 `json:"exposure_units,omitempty"`
	ProductId string `json:"product_id"`
	RiskAdjustmentFactor float64 `json:"risk_adjustment_factor,omitempty"`
}
