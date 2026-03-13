package models


// PolicyRuleDeletionResponse represents a policy_rule_deletion_response
type PolicyRuleDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
