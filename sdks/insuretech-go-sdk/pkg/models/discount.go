package models


// Discount represents a discount
type Discount struct {
	Amount *Money `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	DiscountId string `json:"discount_id,omitempty"`
	DiscountType string `json:"discount_type,omitempty"`
	Name string `json:"name,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}
