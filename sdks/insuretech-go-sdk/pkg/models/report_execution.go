package models

import (
	"time"
)

// ReportExecution represents a report_execution
type ReportExecution struct {
	ReportDefinitionId string `json:"report_definition_id"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	StartedAt time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	ReportScheduleId string `json:"report_schedule_id,omitempty"`
	Parameters string `json:"parameters,omitempty"`
	Status interface{} `json:"status"`
	FileUrl string `json:"file_url,omitempty"`
	FileFormat string `json:"file_format,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ExecutedBy string `json:"executed_by,omitempty"`
	Id string `json:"id"`
}
