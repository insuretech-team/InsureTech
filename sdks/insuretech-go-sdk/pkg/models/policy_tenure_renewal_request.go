package models


// PolicyTenureRenewalRequest represents a policy_tenure_renewal_request
type PolicyTenureRenewalRequest struct {
	Nominees []*Nominee `json:"nominees,omitempty"`
	PolicyId string `json:"policy_id"`
	TenureMonths int `json:"tenure_months,omitempty"`
	UpdateNominees bool `json:"update_nominees,omitempty"`
}
