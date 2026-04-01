package models


// LifeProductsListingResponse represents a life_products_listing_response
type LifeProductsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Products []*LifeProduct `json:"products,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
