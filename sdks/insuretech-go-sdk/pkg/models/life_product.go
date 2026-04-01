package models

import (
	"time"
)

// LifeProduct represents a life_product
type LifeProduct struct {
	AgeAdditionConfig *AgeAdditionConfig `json:"age_addition_config"`
	BaseRate string `json:"base_rate"`
	BonusConfigJson string `json:"bonus_config_json,omitempty"`
	ConditionMultipliersJson string `json:"condition_multipliers_json"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active"`
	MaxEntryAge int `json:"max_entry_age"`
	MaxPolicyTerm int `json:"max_policy_term"`
	MaxSumAssured string `json:"max_sum_assured"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	MinEntryAge int `json:"min_entry_age"`
	MinPolicyTerm int `json:"min_policy_term"`
	MinSumAssured string `json:"min_sum_assured"`
	ProductCode string `json:"product_code"`
	ProductId string `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductType *LifeProductType `json:"product_type"`
	UpdatedAt time.Time `json:"updated_at"`
}
