package models


// PortalConfigUpdateResponse represents a portal_config_update_response
type PortalConfigUpdateResponse struct {
	Config *PortalConfig `json:"config,omitempty"`
	Error *Error `json:"error,omitempty"`
}
