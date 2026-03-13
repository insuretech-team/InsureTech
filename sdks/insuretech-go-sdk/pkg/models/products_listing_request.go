package models


// ProductsListingRequest represents a products_listing_request
type ProductsListingRequest struct {
	Category *ProductCategory `json:"category"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
