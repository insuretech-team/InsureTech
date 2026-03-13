package models


// ProductsProductsListingResponse represents a products_products_listing_response
type ProductsProductsListingResponse struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Error *Error `json:"error,omitempty"`
	Products []*Product `json:"products,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
