package models


// InsuranceQuotesListingResponse represents a insurance_quotes_listing_response
type InsuranceQuotesListingResponse struct {
	Quotes []*UnderwritingQuote `json:"quotes,omitempty"`
	Total int `json:"total,omitempty"`
}
