package models

import (
	"time"
)

// InvoiceCreationRequest represents a invoice_creation_request
type InvoiceCreationRequest struct {
	Amount *Money `json:"amount,omitempty"`
	BusinessId string `json:"business_id"`
	Currency string `json:"currency,omitempty"`
	CustomerId string `json:"customer_id"`
	DueDate time.Time `json:"due_date,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Notes string `json:"notes,omitempty"`
	OrderId string `json:"order_id"`
	OrganisationId string `json:"organisation_id"`
	PolicyIds []string `json:"policy_ids,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	TenantId string `json:"tenant_id"`
}
