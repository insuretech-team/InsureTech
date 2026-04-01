package models

import (
	"time"
)

// LossRatioCalculationResponse represents a loss_ratio_calculation_response
type LossRatioCalculationResponse struct {
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	Errors []string `json:"errors,omitempty"`
	LossRatioId string `json:"loss_ratio_id,omitempty"`
	Result *LossRatioResult `json:"result,omitempty"`
	Success bool `json:"success,omitempty"`
}
