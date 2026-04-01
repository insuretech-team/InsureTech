package models


// ComplianceLogsRetrievalRequest represents a compliance_logs_retrieval_request
type ComplianceLogsRetrievalRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Regulation string `json:"regulation,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	Status string `json:"status,omitempty"`
	Type string `json:"type"`
}
