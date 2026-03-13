package models


// InsuranceBeneficiariesListingResponse represents a insurance_beneficiaries_listing_response
type InsuranceBeneficiariesListingResponse struct {
	Beneficiaries []*Beneficiary `json:"beneficiaries,omitempty"`
	Total int `json:"total,omitempty"`
}
