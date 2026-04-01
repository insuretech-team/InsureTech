package models


// APIKeysListingRequest represents a api_keys_listing_request
type APIKeysListingRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
	OwnerId string `json:"owner_id"`
	OwnerType string `json:"owner_type,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
