package models


// AgeAdditionConfig represents a age_addition_config
type AgeAdditionConfig struct {
	AgeIncrement int `json:"age_increment,omitempty"`
	PriceToAdd string `json:"price_to_add,omitempty"`
	StartAge int `json:"start_age,omitempty"`
}
