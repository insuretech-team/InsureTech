package models


// LossRatioResult represents a loss_ratio_result
type LossRatioResult struct {
	CombinedRatio float64 `json:"combined_ratio,omitempty"`
	ExpenseRatio float64 `json:"expense_ratio,omitempty"`
	Interpretation string `json:"interpretation,omitempty"`
	LossRatio float64 `json:"loss_ratio,omitempty"`
	UnderwritingProfitMargin float64 `json:"underwriting_profit_margin,omitempty"`
}
