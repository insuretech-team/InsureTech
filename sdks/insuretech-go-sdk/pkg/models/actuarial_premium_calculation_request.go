package models


// ActuarialPremiumCalculationRequest represents a actuarial_premium_calculation_request
type ActuarialPremiumCalculationRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	Input *PremiumCalculationInput `json:"input"`
	SaveCalculation bool `json:"save_calculation,omitempty"`
}
