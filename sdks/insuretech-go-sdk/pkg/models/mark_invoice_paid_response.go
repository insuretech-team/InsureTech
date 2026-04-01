package models


// MarkInvoicePaidResponse represents a mark_invoice_paid_response
type MarkInvoicePaidResponse struct {
	InvoiceId string `json:"invoice_id,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
}
