package models


// BusinessRule represents a business_rule
type BusinessRule struct {
	ChildRules []*BusinessRule `json:"child_rules,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ErrorType *ErrorType `json:"error_type,omitempty"`
	Expression string `json:"expression,omitempty"`
	ExpressionType *RuleExpressionType `json:"expression_type,omitempty"`
	IsEnabled bool `json:"is_enabled,omitempty"`
	Priority int `json:"priority,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	RuleId string `json:"rule_id,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
	SuccessEvent string `json:"success_event,omitempty"`
}
