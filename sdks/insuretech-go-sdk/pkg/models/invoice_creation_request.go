package models

import (
	"time"
)

// InvoiceCreationRequest represents a invoice_creation_request
type InvoiceCreationRequest struct {
	CustomerId string `json:"customer_id"`
	BusinessId string `json:"business_id"`
	TenantId string `json:"tenant_id"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	PolicyIds []string `json:"policy_ids,omitempty"`
	OrganisationId string `json:"organisation_id"`
	Amount *Money `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
	Notes string `json:"notes,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	OrderId string `json:"order_id"`
	PurchaseOrderId string `json:"purchase_order_id"`
}
