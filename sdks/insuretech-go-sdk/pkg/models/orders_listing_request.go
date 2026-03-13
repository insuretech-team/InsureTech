package models

import (
	"time"
)

// OrdersListingRequest represents a orders_listing_request
type OrdersListingRequest struct {
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate time.Time `json:"end_date,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	CustomerId string `json:"customer_id"`
	Status *OrderStatus `json:"status,omitempty"`
}
