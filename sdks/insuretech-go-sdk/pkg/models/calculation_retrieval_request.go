package models


// CalculationRetrievalRequest represents a calculation_retrieval_request
type CalculationRetrievalRequest struct {
	CalculationId string `json:"calculation_id"`
	CalculationReference string `json:"calculation_reference,omitempty"`
}
