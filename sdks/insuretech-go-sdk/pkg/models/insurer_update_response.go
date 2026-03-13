package models


// InsurerUpdateResponse represents a insurer_update_response
type InsurerUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
