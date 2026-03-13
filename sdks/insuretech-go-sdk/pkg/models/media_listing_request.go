package models


// MediaListingRequest represents a media_listing_request
type MediaListingRequest struct {
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	MediaType string `json:"media_type,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
