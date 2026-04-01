package models


// UserDocumentUploadRequest represents a user_document_upload_request
type UserDocumentUploadRequest struct {
	DocumentTypeId string `json:"document_type_id"`
	FileUrl string `json:"file_url,omitempty"`
	PolicyId string `json:"policy_id"`
	UserId string `json:"user_id"`
}
