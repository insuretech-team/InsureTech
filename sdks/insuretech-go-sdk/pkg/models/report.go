package models

import (
	"time"
)

// Report represents a report
type Report struct {
	CreatedAt time.Time `json:"created_at"`
	GeneratedBy string `json:"generated_by,omitempty"`
	PeriodEnd time.Time `json:"period_end"`
	PeriodStart time.Time `json:"period_start"`
	ReportData string `json:"report_data"`
	ReportId string `json:"report_id"`
	ReportName string `json:"report_name"`
	ReportUrl string `json:"report_url,omitempty"`
	Status interface{} `json:"status"`
	Type *ReportType `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
