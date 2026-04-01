package models


// MediaListingRequest represents a media_listing_request
type MediaListingRequest struct {
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	MediaType string `json:"media_type,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
}
