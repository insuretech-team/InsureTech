package models

import (
	"time"
)

// QuoteComparison represents a quote_comparison
type QuoteComparison struct {
	CoverageCount int `json:"coverage_count,omitempty"`
	Coverages []*CoverageComparison `json:"coverages,omitempty"`
	DiscountCount int `json:"discount_count,omitempty"`
	Discounts []*DiscountComparison `json:"discounts,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	Status *QuotingQuoteStatus `json:"status,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
