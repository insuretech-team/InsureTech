package models


// DocumentGenerationRequest represents a document_generation_request
type DocumentGenerationRequest struct {
	Alt string `json:"alt,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	IncludeQrCode bool `json:"include_qr_code,omitempty"`
	OutputFormat string `json:"output_format,omitempty"`
	TemplateId string `json:"template_id"`
}
