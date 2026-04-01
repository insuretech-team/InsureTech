package models


// PolicyRuleCreationRequest represents a policy_rule_creation_request
type PolicyRuleCreationRequest struct {
	Action string `json:"action"`
	Condition string `json:"condition,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Description string `json:"description,omitempty"`
	Domain string `json:"domain,omitempty"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	Object string `json:"object,omitempty"`
	Subject string `json:"subject,omitempty"`
}
