package models


// InvoiceIssuanceRequest represents a invoice_issuance_request
type InvoiceIssuanceRequest struct {
	InvoiceId string `json:"invoice_id"`
	IssuedBy string `json:"issued_by,omitempty"`
}
