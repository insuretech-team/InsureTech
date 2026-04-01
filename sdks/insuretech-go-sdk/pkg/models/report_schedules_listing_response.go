package models


// ReportSchedulesListingResponse represents a report_schedules_listing_response
type ReportSchedulesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	ReportSchedules []*ReportSchedule `json:"report_schedules,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
