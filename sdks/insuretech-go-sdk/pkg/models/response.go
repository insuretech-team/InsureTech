package models


// Response represents a response
type Response struct {
	ErrorCode string `json:"error_code,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Success bool `json:"success,omitempty"`
}
