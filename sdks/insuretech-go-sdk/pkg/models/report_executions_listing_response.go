package models


// ReportExecutionsListingResponse represents a report_executions_listing_response
type ReportExecutionsListingResponse struct {
	ReportExecutions []*ReportExecution `json:"report_executions,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
