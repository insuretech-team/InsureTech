package models


// ScheduleReportResponse represents a schedule_report_response
type ScheduleReportResponse struct {
	NextRunAt string `json:"next_run_at,omitempty"`
	ReportId string `json:"report_id,omitempty"`
}
