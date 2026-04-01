package models

import (
	"time"
)

// LossRatioCalculation represents a loss_ratio_calculation
type LossRatioCalculation struct {
	BreakdownJson string `json:"breakdown_json,omitempty"`
	CombinedRatio float64 `json:"combined_ratio,omitempty"`
	ComparisonPeriodId string `json:"comparison_period_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DevelopmentFactorsJson string `json:"development_factors_json,omitempty"`
	EarnedPremium *Money `json:"earned_premium,omitempty"`
	ExpenseRatio float64 `json:"expense_ratio,omitempty"`
	IncurredLosses *Money `json:"incurred_losses,omitempty"`
	LineOfBusiness string `json:"line_of_business,omitempty"`
	LossAdjustmentExpenses *Money `json:"loss_adjustment_expenses,omitempty"`
	LossRatio float64 `json:"loss_ratio,omitempty"`
	LossRatioId string `json:"loss_ratio_id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PeriodEnd time.Time `json:"period_end"`
	PeriodStart time.Time `json:"period_start"`
	ProductId string `json:"product_id,omitempty"`
	TargetLossRatio float64 `json:"target_loss_ratio,omitempty"`
	TotalIncurred *Money `json:"total_incurred,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	WrittenPremium *Money `json:"written_premium,omitempty"`
}
