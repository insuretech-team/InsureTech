package models


// AuditLogCreationResponse represents a audit_log_creation_response
type AuditLogCreationResponse struct {
	AuditLogId string `json:"audit_log_id,omitempty"`
	Error *Error `json:"error,omitempty"`
}
