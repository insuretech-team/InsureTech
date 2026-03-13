package models


// UserDocumentRetrievalResponse represents a user_document_retrieval_response
type UserDocumentRetrievalResponse struct {
	Document *UserDocument `json:"document,omitempty"`
	Error *Error `json:"error,omitempty"`
}
