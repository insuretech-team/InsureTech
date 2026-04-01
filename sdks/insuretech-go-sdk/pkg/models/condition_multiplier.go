package models


// ConditionMultiplier represents a condition_multiplier
type ConditionMultiplier struct {
	ConditionCode string `json:"condition_code,omitempty"`
	ConditionName string `json:"condition_name,omitempty"`
	Description string `json:"description,omitempty"`
	Multiplier float64 `json:"multiplier,omitempty"`
	Severity string `json:"severity,omitempty"`
}
