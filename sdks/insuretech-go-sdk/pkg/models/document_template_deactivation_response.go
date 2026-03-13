package models


// DocumentTemplateDeactivationResponse represents a document_template_deactivation_response
type DocumentTemplateDeactivationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
