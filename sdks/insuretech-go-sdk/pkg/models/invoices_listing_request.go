package models

import (
	"time"
)

// InvoicesListingRequest represents a invoices_listing_request
type InvoicesListingRequest struct {
	PageToken string `json:"page_token,omitempty"`
	CustomerId string `json:"customer_id"`
	EndDate time.Time `json:"end_date,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	OrganisationId string `json:"organisation_id"`
	OrderId string `json:"order_id"`
	PurchaseOrderId string `json:"purchase_order_id"`
	Status *InvoiceStatus `json:"status,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
}
