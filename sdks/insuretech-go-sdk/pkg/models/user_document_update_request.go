package models


// UserDocumentUpdateRequest represents a user_document_update_request
type UserDocumentUpdateRequest struct {
	DocumentTypeId string `json:"document_type_id"`
	FileUrl string `json:"file_url,omitempty"`
	PolicyId string `json:"policy_id"`
	UserDocumentId string `json:"user_document_id"`
}
