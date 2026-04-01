package models


// RatingFormulaRetrievalResponse represents a rating_formula_retrieval_response
type RatingFormulaRetrievalResponse struct {
	Formula *RatingFormula `json:"formula,omitempty"`
	Found bool `json:"found,omitempty"`
}
