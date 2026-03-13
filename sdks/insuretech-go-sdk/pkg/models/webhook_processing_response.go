package models


// WebhookProcessingResponse represents a webhook_processing_response
type WebhookProcessingResponse struct {
	Success bool `json:"success,omitempty"`
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
