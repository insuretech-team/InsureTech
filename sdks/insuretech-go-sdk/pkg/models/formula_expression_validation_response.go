package models


// FormulaExpressionValidationResponse represents a formula_expression_validation_response
type FormulaExpressionValidationResponse struct {
	Errors []string `json:"errors,omitempty"`
	IsValid bool `json:"is_valid,omitempty"`
	ParsedVariables []string `json:"parsed_variables,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
