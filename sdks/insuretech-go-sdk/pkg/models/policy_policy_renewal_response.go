package models


// PolicyPolicyRenewalResponse represents a policy_policy_renewal_response
type PolicyPolicyRenewalResponse struct {
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	NewPolicyId string `json:"new_policy_id,omitempty"`
	NewPolicyNumber string `json:"new_policy_number,omitempty"`
}
