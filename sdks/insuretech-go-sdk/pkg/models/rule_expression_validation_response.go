package models


// RuleExpressionValidationResponse represents a rule_expression_validation_response
type RuleExpressionValidationResponse struct {
	ErrorMessage string `json:"error_message,omitempty"`
	IsValid bool `json:"is_valid,omitempty"`
	TestPassed bool `json:"test_passed,omitempty"`
}
