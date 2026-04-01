package models


// PolicyPolicyCreationRequest represents a policy_policy_creation_request
type PolicyPolicyCreationRequest struct {
	AgentId string `json:"agent_id"`
	Applicant *Applicant `json:"applicant,omitempty"`
	CustomerId string `json:"customer_id"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	PartnerId string `json:"partner_id"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id"`
	Riders []*PolicyRider `json:"riders,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
}
