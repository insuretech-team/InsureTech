package models

import (
	"time"
)

// PaymentsListingRequest represents a payments_listing_request
type PaymentsListingRequest struct {
	EndDate time.Time `json:"end_date,omitempty"`
	UserId string `json:"user_id"`
	PolicyId string `json:"policy_id"`
	Status string `json:"status,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
}
