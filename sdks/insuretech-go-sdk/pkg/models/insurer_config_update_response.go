package models


// InsurerConfigUpdateResponse represents a insurer_config_update_response
type InsurerConfigUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
