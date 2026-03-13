package models


// ErrorInfo represents a error_info
type ErrorInfo struct {
	Code int `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}
