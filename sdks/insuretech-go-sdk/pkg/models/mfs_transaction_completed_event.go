package models

import (
	"time"
)

// MFSTransactionCompletedEvent represents a mfs_transaction_completed_event
type MFSTransactionCompletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	ProviderTransactionId string `json:"provider_transaction_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
