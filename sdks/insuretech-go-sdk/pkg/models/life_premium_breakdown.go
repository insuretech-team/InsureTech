package models


// LifePremiumBreakdown represents a life_premium_breakdown
type LifePremiumBreakdown struct {
	Amount string `json:"amount,omitempty"`
	Component string `json:"component,omitempty"`
	Description string `json:"description,omitempty"`
	IsDiscount bool `json:"is_discount,omitempty"`
}
