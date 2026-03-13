package models


// ErrorResponse represents a error_response
type ErrorResponse struct {
	Code *WebrtcErrorCode `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}
