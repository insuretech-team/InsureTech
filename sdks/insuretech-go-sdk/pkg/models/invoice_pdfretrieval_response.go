package models

import (
	"time"
)

// InvoicePDFRetrievalResponse represents a invoice_pdfretrieval_response
type InvoicePDFRetrievalResponse struct {
	InvoicePdfUrl string `json:"invoice_pdf_url,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	Error *Error `json:"error,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
}
