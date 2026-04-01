package models


// AuditTrailRetrievalRequest represents a audit_trail_retrieval_request
type AuditTrailRetrievalRequest struct {
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}
