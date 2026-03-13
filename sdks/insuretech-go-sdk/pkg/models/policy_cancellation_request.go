package models


// PolicyCancellationRequest represents a policy_cancellation_request
type PolicyCancellationRequest struct {
	PolicyId string `json:"policy_id"`
	Reason string `json:"reason,omitempty"`
}
