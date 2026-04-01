package models

import (
	"time"
)

// GracePeriod represents a grace_period
type GracePeriod struct {
	AuditInfo interface{} `json:"audit_info"`
	CoverageActive bool `json:"coverage_active,omitempty"`
	DaysRemaining int `json:"days_remaining,omitempty"`
	EndDate time.Time `json:"end_date"`
	Id string `json:"id"`
	PolicyId string `json:"policy_id"`
	RevivedAt time.Time `json:"revived_at,omitempty"`
	StartDate time.Time `json:"start_date"`
	Status interface{} `json:"status"`
}
