package models

import (
	"time"
)

// MFSTransaction represents a mfs_transaction
type MFSTransaction struct {
	Amount *Money `json:"amount,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CustomerMsisdn string `json:"customer_msisdn"`
	ErrorMessage string `json:"error_message,omitempty"`
	Id string `json:"id"`
	MfsIntegrationId string `json:"mfs_integration_id"`
	PaymentId string `json:"payment_id,omitempty"`
	Provider string `json:"provider"`
	ProviderTransactionId string `json:"provider_transaction_id,omitempty"`
	RequestPayload string `json:"request_payload,omitempty"`
	ResponsePayload string `json:"response_payload,omitempty"`
	Status interface{} `json:"status"`
	TransactionId string `json:"transaction_id"`
	Type *TransactionType `json:"type"`
}
