package models


// InsurerInsurerConfigUpdateResponse represents a insurer_insurer_config_update_response
type InsurerInsurerConfigUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
