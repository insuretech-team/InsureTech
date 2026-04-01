package models


// PolicyPolicyUpdateRequest represents a policy_policy_update_request
type PolicyPolicyUpdateRequest struct {
	Address string `json:"address,omitempty"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	PolicyId string `json:"policy_id"`
}
