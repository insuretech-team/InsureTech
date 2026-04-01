package models


// ErrorResponse represents a error_response
type ErrorResponse struct {
	Code *WebrtcErrorCode `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}
