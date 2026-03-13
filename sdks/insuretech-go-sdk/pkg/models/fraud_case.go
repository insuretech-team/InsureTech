package models

import (
	"time"
)

// FraudCase represents a fraud_case
type FraudCase struct {
	FraudAlertId string `json:"fraud_alert_id"`
	Priority interface{} `json:"priority"`
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	Status interface{} `json:"status"`
	ClosedAt time.Time `json:"closed_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Id string `json:"id"`
	Evidence string `json:"evidence,omitempty"`
	Outcome *CaseOutcome `json:"outcome,omitempty"`
	InvestigatorId string `json:"investigator_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CaseNumber string `json:"case_number"`
}
