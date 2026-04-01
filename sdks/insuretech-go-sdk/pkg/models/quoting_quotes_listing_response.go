package models


// QuotingQuotesListingResponse represents a quoting_quotes_listing_response
type QuotingQuotesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Quotes []*QuotingQuote `json:"quotes,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
