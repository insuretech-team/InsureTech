package models


// InvoicePDFGenerationRequest represents a invoice_pdfgeneration_request
type InvoicePDFGenerationRequest struct {
	InvoiceId string `json:"invoice_id"`
	RequestedBy string `json:"requested_by,omitempty"`
}
