package models


// PolicyEvaluationResponse represents a policy_evaluation_response
type PolicyEvaluationResponse struct {
	MatchedPolicies []string `json:"matched_policies,omitempty"`
	Error *Error `json:"error,omitempty"`
	Allowed bool `json:"allowed,omitempty"`
	Reason string `json:"reason,omitempty"`
}
