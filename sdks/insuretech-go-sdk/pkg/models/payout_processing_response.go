package models


// PayoutProcessingResponse represents a payout_processing_response
type PayoutProcessingResponse struct {
	Message string `json:"message,omitempty"`
	PaidAt string `json:"paid_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
