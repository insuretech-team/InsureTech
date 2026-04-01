package models


// AuditTrailRetrievalResponse represents a audit_trail_retrieval_response
type AuditTrailRetrievalResponse struct {
	AuditLogs []*AuditLog `json:"audit_logs,omitempty"`
}
