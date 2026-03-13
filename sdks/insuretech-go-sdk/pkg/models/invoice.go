package models

import (
	"time"
)

// Invoice represents a invoice
type Invoice struct {
	OrderId string `json:"order_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	Notes string `json:"notes,omitempty"`
	CancelledAt time.Time `json:"cancelled_at,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	CreditNoteId string `json:"credit_note_id,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	OverdueAt time.Time `json:"overdue_at,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CustomerId string `json:"customer_id,omitempty"`
	Currency string `json:"currency,omitempty"`
	IssuedBy string `json:"issued_by,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	InvoicePdfUrl string `json:"invoice_pdf_url,omitempty"`
	PolicyIds []string `json:"policy_ids,omitempty"`
}
