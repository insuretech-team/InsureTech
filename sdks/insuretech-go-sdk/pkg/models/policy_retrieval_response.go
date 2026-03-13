package models


// PolicyRetrievalResponse represents a policy_retrieval_response
type PolicyRetrievalResponse struct {
	Policy *Policy `json:"policy,omitempty"`
	Error *Error `json:"error,omitempty"`
}
