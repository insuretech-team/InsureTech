package models


// ProductsPremiumBreakdown represents a products_premium_breakdown
type ProductsPremiumBreakdown struct {
	Amount *Money `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	Item string `json:"item,omitempty"`
}
