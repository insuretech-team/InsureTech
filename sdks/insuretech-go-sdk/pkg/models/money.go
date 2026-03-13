package models


// Money represents a money
type Money struct {
	Amount string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
	DecimalAmount float64 `json:"decimal_amount,omitempty"`
}
