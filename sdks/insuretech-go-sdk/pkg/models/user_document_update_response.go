package models


// UserDocumentUpdateResponse represents a user_document_update_response
type UserDocumentUpdateResponse struct {
	Document *UserDocument `json:"document,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
