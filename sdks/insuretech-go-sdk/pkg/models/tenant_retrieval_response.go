package models


// TenantRetrievalResponse represents a tenant_retrieval_response
type TenantRetrievalResponse struct {
	Tenant *Tenant `json:"tenant,omitempty"`
	Error *Error `json:"error,omitempty"`
}
