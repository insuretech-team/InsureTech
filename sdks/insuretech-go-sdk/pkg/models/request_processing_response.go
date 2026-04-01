package models


// RequestProcessingResponse represents a request_processing_response
type RequestProcessingResponse struct {
	JobId string `json:"job_id,omitempty"`
	Status string `json:"status,omitempty"`
}
