package models


// InsurerProductsListingResponse represents a insurer_products_listing_response
type InsurerProductsListingResponse struct {
	Error *Error `json:"error,omitempty"`
	InsurerProducts []*InsurerProduct `json:"insurer_products,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
