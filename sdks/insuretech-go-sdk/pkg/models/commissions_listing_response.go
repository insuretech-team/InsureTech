package models


// CommissionsListingResponse represents a commissions_listing_response
type CommissionsListingResponse struct {
	TotalAmount *Money `json:"total_amount,omitempty"`
	Error *Error `json:"error,omitempty"`
	Commissions []*Commission `json:"commissions,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
