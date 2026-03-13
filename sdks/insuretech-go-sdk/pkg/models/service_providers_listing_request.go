package models


// ServiceProvidersListingRequest represents a service_providers_listing_request
type ServiceProvidersListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	ProviderType string `json:"provider_type"`
	City string `json:"city,omitempty"`
}
