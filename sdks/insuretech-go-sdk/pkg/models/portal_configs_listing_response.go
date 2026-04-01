package models


// PortalConfigsListingResponse represents a portal_configs_listing_response
type PortalConfigsListingResponse struct {
	Configs []*PortalConfig `json:"configs,omitempty"`
}
