package models


// DocumentTemplateCreationRequest represents a document_template_creation_request
type DocumentTemplateCreationRequest struct {
	Description string `json:"description,omitempty"`
	Name string `json:"name"`
	OutputFormat string `json:"output_format,omitempty"`
	TemplateContent string `json:"template_content,omitempty"`
	Type string `json:"type"`
	Variables []string `json:"variables,omitempty"`
}
