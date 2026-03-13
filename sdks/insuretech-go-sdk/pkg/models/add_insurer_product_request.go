package models


// AddInsurerProductRequest represents a add_insurer_product_request
type AddInsurerProductRequest struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name"`
	MinSumAssured *Money `json:"min_sum_assured,omitempty"`
	MaxSumAssured *Money `json:"max_sum_assured,omitempty"`
	MinEntryAge int `json:"min_entry_age,omitempty"`
	MaxEntryAge int `json:"max_entry_age,omitempty"`
	InsurerId string `json:"insurer_id"`
	ProductId string `json:"product_id"`
	EffectiveFrom string `json:"effective_from,omitempty"`
}
