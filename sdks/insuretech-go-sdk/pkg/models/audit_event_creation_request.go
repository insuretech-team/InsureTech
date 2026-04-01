package models


// AuditEventCreationRequest represents a audit_event_creation_request
type AuditEventCreationRequest struct {
	Category string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	EventType string `json:"event_type,omitempty"`
	Metadata string `json:"metadata,omitempty"`
	Severity string `json:"severity,omitempty"`
}
