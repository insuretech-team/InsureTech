package models


// RatingFormulaRetrievalRequest represents a rating_formula_retrieval_request
type RatingFormulaRetrievalRequest struct {
	FormulaCode string `json:"formula_code,omitempty"`
	FormulaId string `json:"formula_id"`
}
