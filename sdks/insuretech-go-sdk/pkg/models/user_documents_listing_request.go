package models


// UserDocumentsListingRequest represents a user_documents_listing_request
type UserDocumentsListingRequest struct {
	PageToken string `json:"page_token,omitempty"`
	UserId string `json:"user_id"`
	DocumentTypeId string `json:"document_type_id"`
	PageSize int `json:"page_size,omitempty"`
}
