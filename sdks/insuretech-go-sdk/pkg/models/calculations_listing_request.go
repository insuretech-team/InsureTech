package models

import (
	"time"
)

// CalculationsListingRequest represents a calculations_listing_request
type CalculationsListingRequest struct {
	CalculatedBy string `json:"calculated_by,omitempty"`
	CalculationType *ActuarialCalculationType `json:"calculation_type,omitempty"`
	DateFrom time.Time `json:"date_from,omitempty"`
	DateTo time.Time `json:"date_to,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
