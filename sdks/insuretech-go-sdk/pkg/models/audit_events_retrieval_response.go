package models


// AuditEventsRetrievalResponse represents a audit_events_retrieval_response
type AuditEventsRetrievalResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	AuditEvents []*AuditEvent `json:"audit_events,omitempty"`
}
