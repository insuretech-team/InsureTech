package models


// TenantCreationResponse represents a tenant_creation_response
type TenantCreationResponse struct {
	TenantId string `json:"tenant_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
