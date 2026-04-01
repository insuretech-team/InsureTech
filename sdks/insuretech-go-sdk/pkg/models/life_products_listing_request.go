package models


// LifeProductsListingRequest represents a life_products_listing_request
type LifeProductsListingRequest struct {
	Filter string `json:"filter,omitempty"`
	OnlyActive bool `json:"only_active,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	ProductType *LifeProductType `json:"product_type"`
}
