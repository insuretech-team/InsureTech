package models


// RenewalPolicyRenewalResponse represents a renewal_policy_renewal_response
type RenewalPolicyRenewalResponse struct {
	NewPolicyId string `json:"new_policy_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
