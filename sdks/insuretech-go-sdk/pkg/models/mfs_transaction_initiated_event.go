package models

import (
	"time"
)

// MFSTransactionInitiatedEvent represents a mfs_transaction_initiated_event
type MFSTransactionInitiatedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
}
