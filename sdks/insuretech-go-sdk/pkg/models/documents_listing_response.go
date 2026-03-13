package models


// DocumentsListingResponse represents a documents_listing_response
type DocumentsListingResponse struct {
	Documents []*DocumentGeneration `json:"documents,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
