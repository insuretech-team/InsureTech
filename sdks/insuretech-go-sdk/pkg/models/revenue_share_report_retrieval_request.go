package models


// RevenueShareReportRetrievalRequest represents a revenue_share_report_retrieval_request
type RevenueShareReportRetrievalRequest struct {
	InsurerId string `json:"insurer_id"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
}
