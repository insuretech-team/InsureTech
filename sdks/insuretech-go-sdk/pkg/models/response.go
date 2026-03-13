package models


// Response represents a response
type Response struct {
	Success bool `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
