package models


// ServiceProvidersListingResponse represents a service_providers_listing_response
type ServiceProvidersListingResponse struct {
	Providers []*ServiceProvider `json:"providers,omitempty"`
	Total int `json:"total,omitempty"`
}
