package models


// InvoiceCancellationResponse represents a invoice_cancellation_response
type InvoiceCancellationResponse struct {
	CreditNoteId string `json:"credit_note_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
}
