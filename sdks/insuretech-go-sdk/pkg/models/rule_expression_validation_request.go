package models


// RuleExpressionValidationRequest represents a rule_expression_validation_request
type RuleExpressionValidationRequest struct {
	Expression string `json:"expression"`
	ExpressionType *RuleExpressionType `json:"expression_type,omitempty"`
	TestInputs map[string]interface{} `json:"test_inputs,omitempty"`
}
