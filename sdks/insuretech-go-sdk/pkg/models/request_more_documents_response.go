package models


// RequestMoreDocumentsResponse represents a request_more_documents_response
type RequestMoreDocumentsResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
