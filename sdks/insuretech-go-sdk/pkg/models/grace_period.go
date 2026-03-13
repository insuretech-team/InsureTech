package models

import (
	"time"
)

// GracePeriod represents a grace_period
type GracePeriod struct {
	Status interface{} `json:"status"`
	RevivedAt time.Time `json:"revived_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	StartDate time.Time `json:"start_date"`
	DaysRemaining int `json:"days_remaining,omitempty"`
	CoverageActive bool `json:"coverage_active,omitempty"`
	PolicyId string `json:"policy_id"`
	EndDate time.Time `json:"end_date"`
}
