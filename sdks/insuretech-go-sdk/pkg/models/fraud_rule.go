package models

import (
	"time"
)

// FraudRule represents a fraud_rule
type FraudRule struct {
	Category *RuleCategory `json:"category"`
	Conditions string `json:"conditions"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	FraudRuleId string `json:"fraud_rule_id"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name"`
	RiskLevel *FraudRiskLevel `json:"risk_level"`
	ScoreWeight int `json:"score_weight"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
