package models


// PolicyPolicyCreationResponse represents a policy_policy_creation_response
type PolicyPolicyCreationResponse struct {
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
