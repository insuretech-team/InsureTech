package models


// PolicyPolicyRetrievalResponse represents a policy_policy_retrieval_response
type PolicyPolicyRetrievalResponse struct {
	Policy *Policy `json:"policy,omitempty"`
	Error *Error `json:"error,omitempty"`
}
