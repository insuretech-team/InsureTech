package models


// InvoiceCancellationRequest represents a invoice_cancellation_request
type InvoiceCancellationRequest struct {
	InvoiceId string `json:"invoice_id"`
	Reason string `json:"reason,omitempty"`
	CancelledBy string `json:"cancelled_by,omitempty"`
	IssueCreditNote bool `json:"issue_credit_note,omitempty"`
}
