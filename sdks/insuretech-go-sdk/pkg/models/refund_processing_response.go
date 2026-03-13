package models


// RefundProcessingResponse represents a refund_processing_response
type RefundProcessingResponse struct {
	Message string `json:"message,omitempty"`
	ProcessedAt string `json:"processed_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
