package models


// OrdersListingResponse represents a orders_listing_response
type OrdersListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Orders []*OrderView `json:"orders,omitempty"`
}
