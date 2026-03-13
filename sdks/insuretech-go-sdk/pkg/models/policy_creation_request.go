package models


// PolicyCreationRequest represents a policy_creation_request
type PolicyCreationRequest struct {
	ProductId string `json:"product_id"`
	PartnerId string `json:"partner_id"`
	Applicant *Applicant `json:"applicant,omitempty"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
	CustomerId string `json:"customer_id"`
	AgentId string `json:"agent_id"`
	Riders []*PolicyRider `json:"riders,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
}
