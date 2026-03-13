package models


// BeneficiaryBeneficiariesListingResponse represents a beneficiary_beneficiaries_listing_response
type BeneficiaryBeneficiariesListingResponse struct {
	Beneficiaries []*Beneficiary `json:"beneficiaries,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
