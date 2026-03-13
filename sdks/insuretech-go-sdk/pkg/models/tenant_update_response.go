package models


// TenantUpdateResponse represents a tenant_update_response
type TenantUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
