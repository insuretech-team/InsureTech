package models


// ComplianceLogsRetrievalRequest represents a compliance_logs_retrieval_request
type ComplianceLogsRetrievalRequest struct {
	Type string `json:"type"`
	Regulation string `json:"regulation,omitempty"`
	Status string `json:"status,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
