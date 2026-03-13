package models


// AuditLogsRetrievalRequest represents a audit_logs_retrieval_request
type AuditLogsRetrievalRequest struct {
	PageSize int `json:"page_size,omitempty"`
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Action string `json:"action"`
	UserId string `json:"user_id"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
}
