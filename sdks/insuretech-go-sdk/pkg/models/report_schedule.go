package models

import (
	"time"
)

// ReportSchedule represents a report_schedule
type ReportSchedule struct {
	AuditInfo interface{} `json:"audit_info"`
	CronExpression string `json:"cron_expression,omitempty"`
	Frequency *ScheduleFrequency `json:"frequency"`
	Id string `json:"id"`
	IsActive bool `json:"is_active,omitempty"`
	LastRunAt time.Time `json:"last_run_at,omitempty"`
	Name string `json:"name"`
	NextRunAt time.Time `json:"next_run_at,omitempty"`
	Parameters string `json:"parameters,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
}
