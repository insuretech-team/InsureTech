package models

import (
	"time"
)

// MFSTransactionFailedEvent represents a mfs_transaction_failed_event
type MFSTransactionFailedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MfsTransactionId string `json:"mfs_transaction_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
