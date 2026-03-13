package models


// BatchCheckAccessResponse represents a batch_check_access_response
type BatchCheckAccessResponse struct {
	Results []*AccessCheckResult `json:"results,omitempty"`
	Error *Error `json:"error,omitempty"`
}
