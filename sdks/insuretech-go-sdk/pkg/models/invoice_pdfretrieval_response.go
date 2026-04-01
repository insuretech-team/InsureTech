package models

import (
	"time"
)

// InvoicePDFRetrievalResponse represents a invoice_pdfretrieval_response
type InvoicePDFRetrievalResponse struct {
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoicePdfUrl string `json:"invoice_pdf_url,omitempty"`
}
