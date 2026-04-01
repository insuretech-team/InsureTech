package models


// RatingFormulaDeletionRequest represents a rating_formula_deletion_request
type RatingFormulaDeletionRequest struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	FormulaId string `json:"formula_id"`
	Permanent bool `json:"permanent,omitempty"`
}
