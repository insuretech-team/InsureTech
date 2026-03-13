package models

import (
	"time"
)

// ReportSchedule represents a report_schedule
type ReportSchedule struct {
	Frequency *ScheduleFrequency `json:"frequency"`
	CronExpression string `json:"cron_expression,omitempty"`
	Parameters string `json:"parameters,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	LastRunAt time.Time `json:"last_run_at,omitempty"`
	Name string `json:"name"`
	Recipients []string `json:"recipients,omitempty"`
	NextRunAt time.Time `json:"next_run_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	ReportDefinitionId string `json:"report_definition_id"`
}
