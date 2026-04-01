package models

import (
	"time"
)

// ActuarialPremiumCalculationResponse represents a actuarial_premium_calculation_response
type ActuarialPremiumCalculationResponse struct {
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationId string `json:"calculation_id,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	Errors []string `json:"errors,omitempty"`
	Result *PremiumCalculationResult `json:"result,omitempty"`
	Success bool `json:"success,omitempty"`
}
