package models


// LifeProductCreationRequest represents a life_product_creation_request
type LifeProductCreationRequest struct {
	AgeAdditionConfig *AgeAdditionConfig `json:"age_addition_config,omitempty"`
	BaseRate string `json:"base_rate,omitempty"`
	Bonuses []*BonusConfig `json:"bonuses,omitempty"`
	ConditionMultipliers []*ConditionMultiplier `json:"condition_multipliers,omitempty"`
	Description string `json:"description,omitempty"`
	MaxEntryAge int `json:"max_entry_age,omitempty"`
	MaxPolicyTerm int `json:"max_policy_term,omitempty"`
	MaxSumAssured string `json:"max_sum_assured,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	MinEntryAge int `json:"min_entry_age,omitempty"`
	MinPolicyTerm int `json:"min_policy_term,omitempty"`
	MinSumAssured string `json:"min_sum_assured,omitempty"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name,omitempty"`
	ProductType *LifeProductType `json:"product_type,omitempty"`
}
