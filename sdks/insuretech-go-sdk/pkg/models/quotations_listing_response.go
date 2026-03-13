package models


// QuotationsListingResponse represents a quotations_listing_response
type QuotationsListingResponse struct {
	Quotations []*Quotation `json:"quotations,omitempty"`
	Total int `json:"total,omitempty"`
}
