package models

import (
	"time"
)

// GatewayWebhookHandlingRequest represents a gateway_webhook_handling_request
type GatewayWebhookHandlingRequest struct {
	Headers map[string]interface{} `json:"headers,omitempty"`
	Provider string `json:"provider"`
	RawPayload string `json:"raw_payload,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
}
