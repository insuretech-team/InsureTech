package models


// InsuranceProductsListingRequest represents a insurance_products_listing_request
type InsuranceProductsListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	TenantId string `json:"tenant_id"`
}
