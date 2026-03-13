package models


// InsurersListingResponse represents a insurers_listing_response
type InsurersListingResponse struct {
	Insurers []*Insurer `json:"insurers,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
