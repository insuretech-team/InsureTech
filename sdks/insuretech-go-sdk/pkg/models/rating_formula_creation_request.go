package models

import (
	"time"
)

// RatingFormulaCreationRequest represents a rating_formula_creation_request
type RatingFormulaCreationRequest struct {
	Category *FormulaCategory `json:"category,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Description string `json:"description,omitempty"`
	FormulaCode string `json:"formula_code"`
	FormulaExpression string `json:"formula_expression,omitempty"`
	FormulaName string `json:"formula_name,omitempty"`
	InsuranceType string `json:"insurance_type,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	SortOrder int `json:"sort_order,omitempty"`
	ValidFrom time.Time `json:"valid_from,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
	Variables []*ActuarialVariable `json:"variables,omitempty"`
}
