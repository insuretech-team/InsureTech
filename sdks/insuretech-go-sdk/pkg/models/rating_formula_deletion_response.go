package models


// RatingFormulaDeletionResponse represents a rating_formula_deletion_response
type RatingFormulaDeletionResponse struct {
	Deleted bool `json:"deleted,omitempty"`
	Success bool `json:"success,omitempty"`
}
