package models


// DocumentTemplatesListingRequest represents a document_templates_listing_request
type DocumentTemplatesListingRequest struct {
	Type string `json:"type"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
