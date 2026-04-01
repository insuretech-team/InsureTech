package models


// PolicyRenewalRequest represents a policy_renewal_request
type PolicyRenewalRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PolicyId string `json:"policy_id"`
}
