package models


// QuotingPremiumBreakdown represents a quoting_premium_breakdown
type QuotingPremiumBreakdown struct {
	Amount *Money `json:"amount,omitempty"`
	Category string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	IsDiscount bool `json:"is_discount,omitempty"`
}
