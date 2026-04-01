package models

import (
	"time"
)

// FraudAlert represents a fraud_alert
type FraudAlert struct {
	AlertNumber string `json:"alert_number"`
	AssignedTo string `json:"assigned_to,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Details string `json:"details,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	FraudRuleId string `json:"fraud_rule_id"`
	FraudScore int `json:"fraud_score"`
	Id string `json:"id"`
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
	RiskLevel string `json:"risk_level"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
