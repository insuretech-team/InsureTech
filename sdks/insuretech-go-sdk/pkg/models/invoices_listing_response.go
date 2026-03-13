package models


// InvoicesListingResponse represents a invoices_listing_response
type InvoicesListingResponse struct {
	Invoices []*Invoice `json:"invoices,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
