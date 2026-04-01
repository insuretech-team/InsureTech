package models

import (
	"time"
)

// RatingFormulaUpdateRequest represents a rating_formula_update_request
type RatingFormulaUpdateRequest struct {
	Category *FormulaCategory `json:"category,omitempty"`
	CreateNewVersion bool `json:"create_new_version,omitempty"`
	Description string `json:"description,omitempty"`
	FormulaExpression string `json:"formula_expression,omitempty"`
	FormulaId string `json:"formula_id"`
	FormulaName string `json:"formula_name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	SortOrder int `json:"sort_order,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
	Variables []*ActuarialVariable `json:"variables,omitempty"`
}
