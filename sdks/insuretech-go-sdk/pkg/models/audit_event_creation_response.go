package models


// AuditEventCreationResponse represents a audit_event_creation_response
type AuditEventCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	AuditEventId string `json:"audit_event_id,omitempty"`
}
