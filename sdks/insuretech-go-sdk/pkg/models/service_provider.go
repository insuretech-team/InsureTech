package models

import (
	"time"
)

// ServiceProvider represents a service_provider
type ServiceProvider struct {
	Address string `json:"address,omitempty"`
	City string `json:"city,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	District string `json:"district,omitempty"`
	Email string `json:"email,omitempty"`
	IsNetworkProvider bool `json:"is_network_provider,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	ProviderId string `json:"provider_id,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	ProviderType *ServiceProviderType `json:"provider_type,omitempty"`
	ServicesOffered []string `json:"services_offered,omitempty"`
	SupportedProductCategories []string `json:"supported_product_categories,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
