package models


// AuditLogCreationRequest represents a audit_log_creation_request
type AuditLogCreationRequest struct {
	Action string `json:"action"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	IpAddress string `json:"ip_address,omitempty"`
	NewValues string `json:"new_values,omitempty"`
	OldValues string `json:"old_values,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
