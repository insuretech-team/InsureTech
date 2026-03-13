package models

import (
	"time"
)

// FraudAlert represents a fraud_alert
type FraudAlert struct {
	EntityId string `json:"entity_id"`
	FraudRuleId string `json:"fraud_rule_id"`
	RiskLevel string `json:"risk_level"`
	Details string `json:"details,omitempty"`
	Status interface{} `json:"status"`
	AssignedTo string `json:"assigned_to,omitempty"`
	EntityType string `json:"entity_type"`
	FraudScore int `json:"fraud_score"`
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Id string `json:"id"`
	AlertNumber string `json:"alert_number"`
}
