package models


// AuditLogCreationRequest represents a audit_log_creation_request
type AuditLogCreationRequest struct {
	UserAgent string `json:"user_agent,omitempty"`
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Action string `json:"action"`
	OldValues string `json:"old_values,omitempty"`
	NewValues string `json:"new_values,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
}
