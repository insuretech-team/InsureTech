package models

import (
	"time"
)

// FraudAlertTriggeredEvent represents a fraud_alert_triggered_event
type FraudAlertTriggeredEvent struct {
	AlertNumber string `json:"alert_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FraudAlertId string `json:"fraud_alert_id,omitempty"`
	FraudScore int `json:"fraud_score,omitempty"`
	RiskLevel string `json:"risk_level,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
