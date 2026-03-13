package models


// InvoiceByOrderIdRetrievalResponse represents a invoice_by_order_id_retrieval_response
type InvoiceByOrderIdRetrievalResponse struct {
	Invoice *Invoice `json:"invoice,omitempty"`
	Error *Error `json:"error,omitempty"`
}
