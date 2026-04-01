package models

import (
	"time"
)

// RatingFormula represents a rating_formula
type RatingFormula struct {
	Category *FormulaCategory `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	FormulaCode string `json:"formula_code"`
	FormulaExpression string `json:"formula_expression"`
	FormulaId string `json:"formula_id"`
	FormulaName string `json:"formula_name"`
	InsuranceType string `json:"insurance_type"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	ParentFormulaId string `json:"parent_formula_id,omitempty"`
	SortOrder int `json:"sort_order"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	ValidFrom time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
	VariablesJson string `json:"variables_json"`
	Version int `json:"version"`
}
