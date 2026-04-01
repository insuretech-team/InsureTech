package models


// LossRatioCalculationRetrievalResponse represents a loss_ratio_calculation_retrieval_response
type LossRatioCalculationRetrievalResponse struct {
	Found bool `json:"found,omitempty"`
	LossRatio *LossRatioCalculation `json:"loss_ratio,omitempty"`
}
