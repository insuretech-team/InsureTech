package models


// CalculationRetrievalResponse represents a calculation_retrieval_response
type CalculationRetrievalResponse struct {
	Calculation *ActuarialCalculation `json:"calculation,omitempty"`
	Found bool `json:"found,omitempty"`
}
