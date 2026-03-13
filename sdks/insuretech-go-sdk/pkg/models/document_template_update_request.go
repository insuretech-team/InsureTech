package models


// DocumentTemplateUpdateRequest represents a document_template_update_request
type DocumentTemplateUpdateRequest struct {
	TemplateId string `json:"template_id"`
	Template *DocumentTemplate `json:"template,omitempty"`
}
