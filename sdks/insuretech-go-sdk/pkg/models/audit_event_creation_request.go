package models


// AuditEventCreationRequest represents a audit_event_creation_request
type AuditEventCreationRequest struct {
	Category string `json:"category,omitempty"`
	EventType string `json:"event_type,omitempty"`
	Severity string `json:"severity,omitempty"`
	Description string `json:"description,omitempty"`
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Metadata string `json:"metadata,omitempty"`
}
