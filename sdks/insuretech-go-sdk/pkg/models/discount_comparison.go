package models


// DiscountComparison represents a discount_comparison
type DiscountComparison struct {
	Amount *Money `json:"amount,omitempty"`
	DiscountId string `json:"discount_id,omitempty"`
	Name string `json:"name,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}
