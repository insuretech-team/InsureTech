package models


// PolicyRuleUpdateRequest represents a policy_rule_update_request
type PolicyRuleUpdateRequest struct {
	PolicyId string `json:"policy_id"`
	Action string `json:"action"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	Condition string `json:"condition,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
