package models


// BusinessRuleResult represents a business_rule_result
type BusinessRuleResult struct {
	ChildResults []string `json:"child_results,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ExecutionTimeMs int `json:"execution_time_ms,omitempty"`
	IsSuccess bool `json:"is_success,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
	SuccessEvent string `json:"success_event,omitempty"`
}
