package models


// DocumentTemplateUpdateRequest represents a document_template_update_request
type DocumentTemplateUpdateRequest struct {
	Template *DocumentTemplate `json:"template,omitempty"`
	TemplateId string `json:"template_id"`
}
