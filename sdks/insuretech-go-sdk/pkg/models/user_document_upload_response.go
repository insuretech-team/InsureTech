package models


// UserDocumentUploadResponse represents a user_document_upload_response
type UserDocumentUploadResponse struct {
	Document *UserDocument `json:"document,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
