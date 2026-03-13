package models


// PolicyRuleUpdateResponse represents a policy_rule_update_response
type PolicyRuleUpdateResponse struct {
	Policy *PolicyRule `json:"policy,omitempty"`
	Error *Error `json:"error,omitempty"`
}
