package models


// AuditTrailRetrievalRequest represents a audit_trail_retrieval_request
type AuditTrailRetrievalRequest struct {
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
}
