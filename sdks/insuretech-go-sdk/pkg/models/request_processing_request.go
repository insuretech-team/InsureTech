package models


// RequestProcessingRequest represents a request_processing_request
type RequestProcessingRequest struct {
	MediaId string `json:"media_id"`
	Options map[string]interface{} `json:"options,omitempty"`
	Priority int `json:"priority,omitempty"`
	ProcessingType string `json:"processing_type,omitempty"`
}
