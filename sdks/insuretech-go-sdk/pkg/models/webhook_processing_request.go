package models


// WebhookProcessingRequest represents a webhook_processing_request
type WebhookProcessingRequest struct {
	Headers map[string]interface{} `json:"headers,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Provider string `json:"provider"`
}
