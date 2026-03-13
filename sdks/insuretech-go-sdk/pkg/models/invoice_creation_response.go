package models


// InvoiceCreationResponse represents a invoice_creation_response
type InvoiceCreationResponse struct {
	Invoice *Invoice `json:"invoice,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
