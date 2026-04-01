package models

import (
	"time"
)

// LossRatioCalculationsListingRequest represents a loss_ratio_calculations_listing_request
type LossRatioCalculationsListingRequest struct {
	LineOfBusiness string `json:"line_of_business,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	PeriodEnd time.Time `json:"period_end,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	ProductId string `json:"product_id"`
}
