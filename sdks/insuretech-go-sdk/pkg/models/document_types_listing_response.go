package models


// DocumentTypesListingResponse represents a document_types_listing_response
type DocumentTypesListingResponse struct {
	Types []*AuthnDocumentType `json:"types,omitempty"`
	Error *Error `json:"error,omitempty"`
}
