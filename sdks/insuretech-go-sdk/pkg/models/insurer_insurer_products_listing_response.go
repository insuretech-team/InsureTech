package models


// InsurerInsurerProductsListingResponse represents a insurer_insurer_products_listing_response
type InsurerInsurerProductsListingResponse struct {
	InsurerProducts []*InsurerProduct `json:"insurer_products,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
