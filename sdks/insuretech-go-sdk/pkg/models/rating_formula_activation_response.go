package models


// RatingFormulaActivationResponse represents a rating_formula_activation_response
type RatingFormulaActivationResponse struct {
	Formula *RatingFormula `json:"formula,omitempty"`
	Success bool `json:"success,omitempty"`
}
