package models


// InsurerInsurerUpdateResponse represents a insurer_insurer_update_response
type InsurerInsurerUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
