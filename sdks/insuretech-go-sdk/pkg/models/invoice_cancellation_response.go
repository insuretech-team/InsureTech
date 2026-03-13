package models


// InvoiceCancellationResponse represents a invoice_cancellation_response
type InvoiceCancellationResponse struct {
	Status *InvoiceStatus `json:"status,omitempty"`
	CreditNoteId string `json:"credit_note_id,omitempty"`
	Error *Error `json:"error,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
}
