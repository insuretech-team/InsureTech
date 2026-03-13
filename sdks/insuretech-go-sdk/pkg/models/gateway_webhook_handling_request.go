package models

import (
	"time"
)

// GatewayWebhookHandlingRequest represents a gateway_webhook_handling_request
type GatewayWebhookHandlingRequest struct {
	Provider string `json:"provider"`
	Headers map[string]interface{} `json:"headers,omitempty"`
	RawPayload string `json:"raw_payload,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
}
