package models


// InvoiceIssuanceResponse represents a invoice_issuance_response
type InvoiceIssuanceResponse struct {
	InvoiceId string `json:"invoice_id,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
	Error *Error `json:"error,omitempty"`
}
