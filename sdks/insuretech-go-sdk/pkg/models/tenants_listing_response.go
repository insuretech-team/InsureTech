package models


// TenantsListingResponse represents a tenants_listing_response
type TenantsListingResponse struct {
	Tenants []*Tenant `json:"tenants,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
