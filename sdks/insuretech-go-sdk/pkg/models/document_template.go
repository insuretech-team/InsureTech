package models


// DocumentTemplate represents a document_template
type DocumentTemplate struct {
	AuditInfo interface{} `json:"audit_info"`
	Description string `json:"description,omitempty"`
	Id string `json:"id"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name"`
	OutputFormat *OutputFormat `json:"output_format"`
	TemplateContent string `json:"template_content"`
	Type *DocumentDocumentType `json:"type"`
	Variables string `json:"variables,omitempty"`
	Version int `json:"version"`
}
