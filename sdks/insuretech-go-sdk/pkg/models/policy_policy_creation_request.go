package models


// PolicyPolicyCreationRequest represents a policy_policy_creation_request
type PolicyPolicyCreationRequest struct {
	ProductId string `json:"product_id"`
	Applicant *Applicant `json:"applicant,omitempty"`
	Riders []*PolicyRider `json:"riders,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
	CustomerId string `json:"customer_id"`
	PartnerId string `json:"partner_id"`
	AgentId string `json:"agent_id"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
}
