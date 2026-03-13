package models

import (
	"time"
)

// ProductsRider represents a products_rider
type ProductsRider struct {
	ProductId string `json:"product_id,omitempty"`
	Description string `json:"description,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	IsMandatory bool `json:"is_mandatory,omitempty"`
	PremiumCurrency string `json:"premium_currency,omitempty"`
	RiderName string `json:"rider_name,omitempty"`
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CoverageCurrency string `json:"coverage_currency,omitempty"`
	RiderId string `json:"rider_id,omitempty"`
}
