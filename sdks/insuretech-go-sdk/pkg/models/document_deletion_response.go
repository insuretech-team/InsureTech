package models


// DocumentDeletionResponse represents a document_deletion_response
type DocumentDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
