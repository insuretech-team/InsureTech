package models


// AuditLogsRetrievalRequest represents a audit_logs_retrieval_request
type AuditLogsRetrievalRequest struct {
	Action string `json:"action"`
	EndDate string `json:"end_date,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	UserId string `json:"user_id"`
}
