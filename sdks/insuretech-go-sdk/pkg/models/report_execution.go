package models

import (
	"time"
)

// ReportExecution represents a report_execution
type ReportExecution struct {
	AuditInfo interface{} `json:"audit_info"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ExecutedBy string `json:"executed_by,omitempty"`
	FileFormat string `json:"file_format,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	Id string `json:"id"`
	Parameters string `json:"parameters,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
	ReportScheduleId string `json:"report_schedule_id,omitempty"`
	RowCount int `json:"row_count,omitempty"`
	StartedAt time.Time `json:"started_at"`
	Status interface{} `json:"status"`
}
