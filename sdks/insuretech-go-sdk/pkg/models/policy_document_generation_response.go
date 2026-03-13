package models


// PolicyDocumentGenerationResponse represents a policy_document_generation_response
type PolicyDocumentGenerationResponse struct {
	DocumentUrl string `json:"document_url,omitempty"`
	QrCode string `json:"qr_code,omitempty"`
	Error *Error `json:"error,omitempty"`
}
