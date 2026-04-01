package models

import (
	"time"
)

// PolicyRider represents a policy_rider
type PolicyRider struct {
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CoverageCurrency string `json:"coverage_currency,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	PremiumCurrency string `json:"premium_currency,omitempty"`
	RiderId string `json:"rider_id,omitempty"`
	RiderName string `json:"rider_name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
