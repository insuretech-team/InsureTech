package models


// RatingFormulaCreationResponse represents a rating_formula_creation_response
type RatingFormulaCreationResponse struct {
	Errors []string `json:"errors,omitempty"`
	Formula *RatingFormula `json:"formula,omitempty"`
	FormulaId string `json:"formula_id,omitempty"`
	Success bool `json:"success,omitempty"`
}
