package models


// DocumentTemplateCreationResponse represents a document_template_creation_response
type DocumentTemplateCreationResponse struct {
	TemplateId string `json:"template_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
