package models


// PolicyIssuanceResponse represents a policy_issuance_response
type PolicyIssuanceResponse struct {
	Policy *Policy `json:"policy,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
