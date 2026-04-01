package models


// DocumentTemplateDeactivationRequest represents a document_template_deactivation_request
type DocumentTemplateDeactivationRequest struct {
	Reason string `json:"reason,omitempty"`
	TemplateId string `json:"template_id"`
}
