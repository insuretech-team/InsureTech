package models


// RatingFormulaUpdateResponse represents a rating_formula_update_response
type RatingFormulaUpdateResponse struct {
	Errors []string `json:"errors,omitempty"`
	Formula *RatingFormula `json:"formula,omitempty"`
	NewVersionCreated bool `json:"new_version_created,omitempty"`
	Success bool `json:"success,omitempty"`
}
