package models

import (
	"time"
)

// ActuarialReportGenerationResponse represents a actuarial_report_generation_response
type ActuarialReportGenerationResponse struct {
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	ReportId string `json:"report_id,omitempty"`
	ReportUrl string `json:"report_url,omitempty"`
	Success bool `json:"success,omitempty"`
	Summary map[string]interface{} `json:"summary,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
