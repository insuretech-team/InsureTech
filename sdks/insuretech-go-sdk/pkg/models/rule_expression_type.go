package models

// RuleExpressionType represents a rule_expression_type
type RuleExpressionType string

// RuleExpressionType values
const (
	RuleExpressionTypeRULEEXPRESSIONTYPEUNSPECIFIED RuleExpressionType = "RULE_EXPRESSION_TYPE_UNSPECIFIED"
	RuleExpressionTypeRULEEXPRESSIONTYPELAMBDA  = "RULE_EXPRESSION_TYPE_LAMBDA"
	RuleExpressionTypeRULEEXPRESSIONTYPECUSTOM  = "RULE_EXPRESSION_TYPE_CUSTOM"
)
