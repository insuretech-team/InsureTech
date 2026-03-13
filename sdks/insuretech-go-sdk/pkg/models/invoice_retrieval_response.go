package models


// InvoiceRetrievalResponse represents a invoice_retrieval_response
type InvoiceRetrievalResponse struct {
	Invoice *Invoice `json:"invoice,omitempty"`
	Error *Error `json:"error,omitempty"`
}
