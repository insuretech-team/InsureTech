package models


// DocumentTemplateDeactivationRequest represents a document_template_deactivation_request
type DocumentTemplateDeactivationRequest struct {
	TemplateId string `json:"template_id"`
	Reason string `json:"reason,omitempty"`
}
