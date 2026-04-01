package models


// PolicyRuleUpdateRequest represents a policy_rule_update_request
type PolicyRuleUpdateRequest struct {
	Action string `json:"action"`
	Condition string `json:"condition,omitempty"`
	Description string `json:"description,omitempty"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	PolicyId string `json:"policy_id"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
