package models


// InsuranceQuotesListingResponse represents a insurance_quotes_listing_response
type InsuranceQuotesListingResponse struct {
	Quotes []*Quote `json:"quotes,omitempty"`
	Total int `json:"total,omitempty"`
}
