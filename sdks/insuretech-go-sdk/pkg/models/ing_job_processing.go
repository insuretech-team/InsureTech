package models

import (
	"time"
)

// ingJobProcessing represents a ing_job_processing
type ingJobProcessing struct {
	MaxRetries int `json:"max_retries"`
	StartedAt time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ResultData string `json:"result_data,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	MediaId string `json:"media_id"`
	Status interface{} `json:"status"`
	Priority int `json:"priority"`
	ProcessingType *ProcessingType `json:"processing_type"`
	RetryCount int `json:"retry_count"`
}
