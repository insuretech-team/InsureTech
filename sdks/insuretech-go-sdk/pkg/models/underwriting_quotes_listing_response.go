package models


// UnderwritingQuotesListingResponse represents a underwriting_quotes_listing_response
type UnderwritingQuotesListingResponse struct {
	Quotes []*UnderwritingQuote `json:"quotes,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
