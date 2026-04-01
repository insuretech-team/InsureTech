package models

import (
	"time"
)

// FraudCase represents a fraud_case
type FraudCase struct {
	CaseNumber string `json:"case_number"`
	ClosedAt time.Time `json:"closed_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Evidence string `json:"evidence,omitempty"`
	FraudAlertId string `json:"fraud_alert_id"`
	Id string `json:"id"`
	InvestigationNotes string `json:"investigation_notes,omitempty"`
	InvestigatorId string `json:"investigator_id,omitempty"`
	Outcome *CaseOutcome `json:"outcome,omitempty"`
	Priority interface{} `json:"priority"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
