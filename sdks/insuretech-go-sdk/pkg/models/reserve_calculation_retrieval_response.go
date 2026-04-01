package models


// ReserveCalculationRetrievalResponse represents a reserve_calculation_retrieval_response
type ReserveCalculationRetrievalResponse struct {
	Found bool `json:"found,omitempty"`
	Reserve *ReserveCalculation `json:"reserve,omitempty"`
}
