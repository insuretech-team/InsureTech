package models


// ApiResponse represents a api_response
type ApiResponse struct {
	Data interface{} `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
	Meta *ResponseMeta `json:"meta,omitempty"`
	Success bool `json:"success"`
}
