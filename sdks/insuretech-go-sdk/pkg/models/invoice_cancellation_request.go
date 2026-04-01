package models


// InvoiceCancellationRequest represents a invoice_cancellation_request
type InvoiceCancellationRequest struct {
	CancelledBy string `json:"cancelled_by,omitempty"`
	InvoiceId string `json:"invoice_id"`
	IssueCreditNote bool `json:"issue_credit_note,omitempty"`
	Reason string `json:"reason,omitempty"`
}
