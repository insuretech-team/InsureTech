package models


// ReportSchedulesListingResponse represents a report_schedules_listing_response
type ReportSchedulesListingResponse struct {
	ReportSchedules []*ReportSchedule `json:"report_schedules,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
