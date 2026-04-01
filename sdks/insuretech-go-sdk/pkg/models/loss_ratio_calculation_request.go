package models

import (
	"time"
)

// LossRatioCalculationRequest represents a loss_ratio_calculation_request
type LossRatioCalculationRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationReference string `json:"calculation_reference,omitempty"`
	Input *LossRatioInput `json:"input,omitempty"`
	LineOfBusiness string `json:"line_of_business,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	ProductId string `json:"product_id"`
	SaveCalculation bool `json:"save_calculation,omitempty"`
}
