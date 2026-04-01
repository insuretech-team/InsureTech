package models

import (
	"time"
)

// MFSWebhook represents a mfs_webhook
type MFSWebhook struct {
	AuditInfo interface{} `json:"audit_info"`
	ErrorMessage string `json:"error_message,omitempty"`
	EventType string `json:"event_type"`
	Headers string `json:"headers,omitempty"`
	Id string `json:"id"`
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	Payload string `json:"payload"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	Provider string `json:"provider"`
	SignatureValid bool `json:"signature_valid,omitempty"`
	Status interface{} `json:"status"`
}
