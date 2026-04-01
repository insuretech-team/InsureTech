package models


// CommissionsListingResponse represents a commissions_listing_response
type CommissionsListingResponse struct {
	Commissions []*Commission `json:"commissions,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
