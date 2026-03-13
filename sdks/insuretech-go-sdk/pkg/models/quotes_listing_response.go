package models


// QuotesListingResponse represents a quotes_listing_response
type QuotesListingResponse struct {
	Quotes []*Quote `json:"quotes,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
