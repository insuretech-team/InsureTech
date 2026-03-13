package models


// DocumentTemplateUpdateResponse represents a document_template_update_response
type DocumentTemplateUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
