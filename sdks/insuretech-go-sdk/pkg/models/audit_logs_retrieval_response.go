package models


// AuditLogsRetrievalResponse represents a audit_logs_retrieval_response
type AuditLogsRetrievalResponse struct {
	AuditLogs []*AuditLog `json:"audit_logs,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
