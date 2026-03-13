package models


// PolicyRuleCreationResponse represents a policy_rule_creation_response
type PolicyRuleCreationResponse struct {
	Policy *PolicyRule `json:"policy,omitempty"`
	Error *Error `json:"error,omitempty"`
}
