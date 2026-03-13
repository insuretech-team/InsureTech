package models


// InsurerInsurersListingResponse represents a insurer_insurers_listing_response
type InsurerInsurersListingResponse struct {
	Insurers []*Insurer `json:"insurers,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
