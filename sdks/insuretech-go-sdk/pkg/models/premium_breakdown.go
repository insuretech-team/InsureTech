package models


// PremiumBreakdown represents a premium_breakdown
type PremiumBreakdown struct {
	Amount *Money `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	Item string `json:"item,omitempty"`
}
