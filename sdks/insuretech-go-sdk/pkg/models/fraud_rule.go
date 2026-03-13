package models

import (
	"time"
)

// FraudRule represents a fraud_rule
type FraudRule struct {
	Category *RuleCategory `json:"category"`
	ScoreWeight int `json:"score_weight"`
	IsActive bool `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	FraudRuleId string `json:"fraud_rule_id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Conditions string `json:"conditions"`
	RiskLevel *FraudRiskLevel `json:"risk_level"`
}
