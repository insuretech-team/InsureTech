package models


// RefundsListingResponse represents a refunds_listing_response
type RefundsListingResponse struct {
	Refunds []*Refund `json:"refunds,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
