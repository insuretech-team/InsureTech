package models


// UserDocumentDeletionResponse represents a user_document_deletion_response
type UserDocumentDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
