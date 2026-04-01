package models


// LifeQuotesListingResponse represents a life_quotes_listing_response
type LifeQuotesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Quotes []*LifeQuote `json:"quotes,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
