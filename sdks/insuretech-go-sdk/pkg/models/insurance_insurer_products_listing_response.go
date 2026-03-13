package models


// InsuranceInsurerProductsListingResponse represents a insurance_insurer_products_listing_response
type InsuranceInsurerProductsListingResponse struct {
	Products []*InsurerProduct `json:"products,omitempty"`
}
