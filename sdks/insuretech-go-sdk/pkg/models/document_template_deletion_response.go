package models


// DocumentTemplateDeletionResponse represents a document_template_deletion_response
type DocumentTemplateDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
