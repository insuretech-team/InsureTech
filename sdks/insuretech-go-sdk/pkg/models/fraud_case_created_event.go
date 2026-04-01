package models

import (
	"time"
)

// FraudCaseCreatedEvent represents a fraud_case_created_event
type FraudCaseCreatedEvent struct {
	CaseNumber string `json:"case_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FraudAlertId string `json:"fraud_alert_id,omitempty"`
	FraudCaseId string `json:"fraud_case_id,omitempty"`
	InvestigatorId string `json:"investigator_id,omitempty"`
	Priority string `json:"priority,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
