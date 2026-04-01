package models


// AuditEventsRetrievalResponse represents a audit_events_retrieval_response
type AuditEventsRetrievalResponse struct {
	AuditEvents []*AuditEvent `json:"audit_events,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
