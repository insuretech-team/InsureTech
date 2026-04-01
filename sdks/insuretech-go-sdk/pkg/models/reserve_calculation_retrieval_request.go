package models


// ReserveCalculationRetrievalRequest represents a reserve_calculation_retrieval_request
type ReserveCalculationRetrievalRequest struct {
	ClaimId string `json:"claim_id"`
	ReserveId string `json:"reserve_id"`
}
