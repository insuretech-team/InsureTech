package models


// RevivePolicyResponse represents a revive_policy_response
type RevivePolicyResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
