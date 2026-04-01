package models

import (
	"time"
)

// LossRatioInput represents a loss_ratio_input
type LossRatioInput struct {
	EarnedPremium float64 `json:"earned_premium,omitempty"`
	IncurredLosses float64 `json:"incurred_losses,omitempty"`
	LineOfBusiness string `json:"line_of_business,omitempty"`
	LossAdjustmentExpenses float64 `json:"loss_adjustment_expenses,omitempty"`
	OperatingExpenses float64 `json:"operating_expenses,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	WrittenPremium float64 `json:"written_premium,omitempty"`
}
