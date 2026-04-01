package models

import (
	"time"
)

// PaymentsListingRequest represents a payments_listing_request
type PaymentsListingRequest struct {
	EndDate time.Time `json:"end_date,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	PolicyId string `json:"policy_id"`
	StartDate time.Time `json:"start_date,omitempty"`
	Status string `json:"status,omitempty"`
	UserId string `json:"user_id"`
}
