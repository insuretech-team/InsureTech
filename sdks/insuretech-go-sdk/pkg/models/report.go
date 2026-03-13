package models

import (
	"time"
)

// Report represents a report
type Report struct {
	PeriodEnd time.Time `json:"period_end"`
	ReportData string `json:"report_data"`
	ReportUrl string `json:"report_url,omitempty"`
	Status interface{} `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ReportId string `json:"report_id"`
	PeriodStart time.Time `json:"period_start"`
	GeneratedBy string `json:"generated_by,omitempty"`
	ReportName string `json:"report_name"`
	Type *ReportType `json:"type"`
}
