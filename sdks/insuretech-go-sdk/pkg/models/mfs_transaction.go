package models

import (
	"time"
)

// MFSTransaction represents a mfs_transaction
type MFSTransaction struct {
	Amount *Money `json:"amount,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Id string `json:"id"`
	MfsIntegrationId string `json:"mfs_integration_id"`
	Provider string `json:"provider"`
	Type *TransactionType `json:"type"`
	CustomerMsisdn string `json:"customer_msisdn"`
	AuditInfo interface{} `json:"audit_info"`
	PaymentId string `json:"payment_id,omitempty"`
	ResponsePayload string `json:"response_payload,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	RequestPayload string `json:"request_payload,omitempty"`
	TransactionId string `json:"transaction_id"`
	ProviderTransactionId string `json:"provider_transaction_id,omitempty"`
	Status interface{} `json:"status"`
}
