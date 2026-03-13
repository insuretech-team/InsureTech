package models


// PortalConfigRetrievalResponse represents a portal_config_retrieval_response
type PortalConfigRetrievalResponse struct {
	Config *PortalConfig `json:"config,omitempty"`
	Error *Error `json:"error,omitempty"`
}
