package models

import (
	"time"
)

// MFSTransactionInitiatedEvent represents a mfs_transaction_initiated_event
type MFSTransactionInitiatedEvent struct {
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}
