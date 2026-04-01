package models


// FormulaExpressionValidationRequest represents a formula_expression_validation_request
type FormulaExpressionValidationRequest struct {
	FormulaExpression string `json:"formula_expression"`
	Variables []*ActuarialVariable `json:"variables,omitempty"`
}
