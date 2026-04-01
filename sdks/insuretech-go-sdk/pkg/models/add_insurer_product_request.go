package models


// AddInsurerProductRequest represents a add_insurer_product_request
type AddInsurerProductRequest struct {
	Code string `json:"code,omitempty"`
	EffectiveFrom string `json:"effective_from,omitempty"`
	InsurerId string `json:"insurer_id"`
	MaxEntryAge int `json:"max_entry_age,omitempty"`
	MaxSumAssured *Money `json:"max_sum_assured,omitempty"`
	MinEntryAge int `json:"min_entry_age,omitempty"`
	MinSumAssured *Money `json:"min_sum_assured,omitempty"`
	Name string `json:"name"`
	ProductId string `json:"product_id"`
}
