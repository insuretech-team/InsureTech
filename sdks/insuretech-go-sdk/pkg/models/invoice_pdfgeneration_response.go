package models


// InvoicePDFGenerationResponse represents a invoice_pdfgeneration_response
type InvoicePDFGenerationResponse struct {
	InvoiceId string `json:"invoice_id,omitempty"`
	JobId string `json:"job_id,omitempty"`
}
