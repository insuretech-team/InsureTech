package models


// PaymentsListingResponse represents a payments_listing_response
type PaymentsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Payments []*Payment `json:"payments,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
