package models


// RevivePolicyRequest represents a revive_policy_request
type RevivePolicyRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PolicyId string `json:"policy_id"`
}
