package models


// DLRStatusUpdateResponse represents a dlrstatus_update_response
type DLRStatusUpdateResponse struct {
	Updated bool `json:"updated,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
