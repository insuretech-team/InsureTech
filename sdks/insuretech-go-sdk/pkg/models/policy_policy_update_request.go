package models


// PolicyPolicyUpdateRequest represents a policy_policy_update_request
type PolicyPolicyUpdateRequest struct {
	PolicyId string `json:"policy_id"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	Address string `json:"address,omitempty"`
}
