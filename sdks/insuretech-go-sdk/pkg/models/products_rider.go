package models

import (
	"time"
)

// ProductsRider represents a products_rider
type ProductsRider struct {
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CoverageCurrency string `json:"coverage_currency,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Description string `json:"description,omitempty"`
	IsMandatory bool `json:"is_mandatory,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	PremiumCurrency string `json:"premium_currency,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	RiderId string `json:"rider_id,omitempty"`
	RiderName string `json:"rider_name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
