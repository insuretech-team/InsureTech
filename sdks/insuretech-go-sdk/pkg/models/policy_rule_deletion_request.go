package models


// PolicyRuleDeletionRequest represents a policy_rule_deletion_request
type PolicyRuleDeletionRequest struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	PolicyId string `json:"policy_id"`
}
