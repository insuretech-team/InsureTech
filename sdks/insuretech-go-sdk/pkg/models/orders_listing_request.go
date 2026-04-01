package models

import (
	"time"
)

// OrdersListingRequest represents a orders_listing_request
type OrdersListingRequest struct {
	CustomerId string `json:"customer_id"`
	EndDate time.Time `json:"end_date,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
}
