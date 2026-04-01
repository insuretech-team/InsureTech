package models

import (
	"time"
)

// ingJobProcessing represents a ing_job_processing
type ingJobProcessing struct {
	AuditInfo interface{} `json:"audit_info"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Id string `json:"id"`
	MaxRetries int `json:"max_retries"`
	MediaId string `json:"media_id"`
	Priority int `json:"priority"`
	ProcessingType *ProcessingType `json:"processing_type"`
	ResultData string `json:"result_data,omitempty"`
	RetryCount int `json:"retry_count"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Status interface{} `json:"status"`
}
