package models


// PolicyPolicyUpdateResponse represents a policy_policy_update_response
type PolicyPolicyUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
