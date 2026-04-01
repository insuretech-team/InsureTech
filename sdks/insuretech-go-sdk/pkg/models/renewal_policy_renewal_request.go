package models


// RenewalPolicyRenewalRequest represents a renewal_policy_renewal_request
type RenewalPolicyRenewalRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PolicyId string `json:"policy_id"`
}
