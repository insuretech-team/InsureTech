package models


// DocumentsListingRequest represents a documents_listing_request
type DocumentsListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Status string `json:"status,omitempty"`
}
