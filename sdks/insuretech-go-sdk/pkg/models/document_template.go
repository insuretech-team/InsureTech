package models


// DocumentTemplate represents a document_template
type DocumentTemplate struct {
	Name string `json:"name"`
	Type *DocumentDocumentType `json:"type"`
	Description string `json:"description,omitempty"`
	TemplateContent string `json:"template_content"`
	Version int `json:"version"`
	IsActive bool `json:"is_active,omitempty"`
	Id string `json:"id"`
	OutputFormat *OutputFormat `json:"output_format"`
	Variables string `json:"variables,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
}
