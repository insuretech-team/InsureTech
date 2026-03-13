package models


// InsuranceInsurersListingResponse represents a insurance_insurers_listing_response
type InsuranceInsurersListingResponse struct {
	Insurers []*Insurer `json:"insurers,omitempty"`
	Total int `json:"total,omitempty"`
}
