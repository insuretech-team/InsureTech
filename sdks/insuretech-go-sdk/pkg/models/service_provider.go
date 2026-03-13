package models

import (
	"time"
)

// ServiceProvider represents a service_provider
type ServiceProvider struct {
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	City string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	ProviderId string `json:"provider_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email string `json:"email,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	ServicesOffered []string `json:"services_offered,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	IsNetworkProvider bool `json:"is_network_provider,omitempty"`
	ProviderType *ServiceProviderType `json:"provider_type,omitempty"`
	Address string `json:"address,omitempty"`
	SupportedProductCategories []string `json:"supported_product_categories,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
