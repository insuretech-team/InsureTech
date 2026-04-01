package models


// PolicyPolicyRenewalResponse represents a policy_policy_renewal_response
type PolicyPolicyRenewalResponse struct {
	NewPolicyId string `json:"new_policy_id,omitempty"`
	NewPolicyNumber string `json:"new_policy_number,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
}
