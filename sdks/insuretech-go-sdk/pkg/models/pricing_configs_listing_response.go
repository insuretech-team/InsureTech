package models


// PricingConfigsListingResponse represents a pricing_configs_listing_response
type PricingConfigsListingResponse struct {
	Configs []*PricingConfig `json:"configs,omitempty"`
}
