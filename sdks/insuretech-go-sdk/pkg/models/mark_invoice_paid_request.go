package models

import (
	"time"
)

// MarkInvoicePaidRequest represents a mark_invoice_paid_request
type MarkInvoicePaidRequest struct {
	InvoiceId string `json:"invoice_id"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	PaymentId string `json:"payment_id"`
}
