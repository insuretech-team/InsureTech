package models


// DocumentTemplatesListingRequest represents a document_templates_listing_request
type DocumentTemplatesListingRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Type string `json:"type"`
}
