package models


// PolicyRuleCreationRequest represents a policy_rule_creation_request
type PolicyRuleCreationRequest struct {
	CreatedBy string `json:"created_by,omitempty"`
	Subject string `json:"subject,omitempty"`
	Domain string `json:"domain,omitempty"`
	Object string `json:"object,omitempty"`
	Action string `json:"action"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	Condition string `json:"condition,omitempty"`
	Description string `json:"description,omitempty"`
}
