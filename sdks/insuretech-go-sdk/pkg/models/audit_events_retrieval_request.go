package models


// AuditEventsRetrievalRequest represents a audit_events_retrieval_request
type AuditEventsRetrievalRequest struct {
	Category string `json:"category"`
	Severity string `json:"severity,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
