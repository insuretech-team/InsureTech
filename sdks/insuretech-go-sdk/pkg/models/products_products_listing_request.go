package models


// ProductsProductsListingRequest represents a products_products_listing_request
type ProductsProductsListingRequest struct {
	Category *ProductCategory `json:"category"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
