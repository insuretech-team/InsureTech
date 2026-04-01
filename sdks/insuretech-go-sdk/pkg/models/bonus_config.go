package models


// BonusConfig represents a bonus_config
type BonusConfig struct {
	BonusCode string `json:"bonus_code,omitempty"`
	BonusName string `json:"bonus_name,omitempty"`
	BonusType string `json:"bonus_type,omitempty"`
	Description string `json:"description,omitempty"`
	FixedAmount string `json:"fixed_amount,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}
