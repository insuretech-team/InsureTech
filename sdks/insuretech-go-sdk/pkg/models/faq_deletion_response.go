package models


// FAQDeletionResponse represents a faq_deletion_response
type FAQDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
