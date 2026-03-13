package models


// ComplianceLogsRetrievalResponse represents a compliance_logs_retrieval_response
type ComplianceLogsRetrievalResponse struct {
	ComplianceLogs []*ComplianceLog `json:"compliance_logs,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
