package models


// InsuranceProductsListingResponse represents a insurance_products_listing_response
type InsuranceProductsListingResponse struct {
	Products []*Product `json:"products,omitempty"`
	Total int `json:"total,omitempty"`
}
