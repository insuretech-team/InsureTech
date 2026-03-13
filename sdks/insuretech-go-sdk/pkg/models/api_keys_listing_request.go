package models


// APIKeysListingRequest represents a api_keys_listing_request
type APIKeysListingRequest struct {
	OwnerType string `json:"owner_type,omitempty"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	OwnerId string `json:"owner_id"`
}
