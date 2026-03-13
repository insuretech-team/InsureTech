package models


// PaymentsListingResponse represents a payments_listing_response
type PaymentsListingResponse struct {
	Payments []*Payment `json:"payments,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
