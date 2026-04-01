package models

import (
	"time"
)

// InvoicesListingRequest represents a invoices_listing_request
type InvoicesListingRequest struct {
	CustomerId string `json:"customer_id"`
	EndDate time.Time `json:"end_date,omitempty"`
	OrderId string `json:"order_id"`
	OrganisationId string `json:"organisation_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id"`
	StartDate time.Time `json:"start_date,omitempty"`
	Status *InvoiceStatus `json:"status,omitempty"`
}
