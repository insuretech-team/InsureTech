package models


// DocumentsListingRequest represents a documents_listing_request
type DocumentsListingRequest struct {
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
}
