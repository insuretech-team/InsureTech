package models


// PolicyTenureRenewalResponse represents a policy_tenure_renewal_response
type PolicyTenureRenewalResponse struct {
	NewPolicyId string `json:"new_policy_id,omitempty"`
	NewPolicyNumber string `json:"new_policy_number,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
}
