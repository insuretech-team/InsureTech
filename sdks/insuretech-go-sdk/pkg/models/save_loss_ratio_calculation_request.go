package models


// SaveLossRatioCalculationRequest represents a save_loss_ratio_calculation_request
type SaveLossRatioCalculationRequest struct {
	LossRatio *LossRatioCalculation `json:"loss_ratio"`
}
