package models

import (
	"time"
)

// Invoice represents a invoice
type Invoice struct {
	Amount *Money `json:"amount,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	CancelledAt time.Time `json:"cancelled_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CreditNoteId string `json:"credit_note_id,omitempty"`
	Currency string `json:"currency,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	InvoicePdfUrl string `json:"invoice_pdf_url,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	IssuedBy string `json:"issued_by,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Notes string `json:"notes,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OverdueAt time.Time `json:"overdue_at,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PolicyIds []string `json:"policy_ids,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
