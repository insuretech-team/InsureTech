package models


// SaveLossRatioCalculationResponse represents a save_loss_ratio_calculation_response
type SaveLossRatioCalculationResponse struct {
	Errors []string `json:"errors,omitempty"`
	LossRatio *LossRatioCalculation `json:"loss_ratio,omitempty"`
	Success bool `json:"success,omitempty"`
}
