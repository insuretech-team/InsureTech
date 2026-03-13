package models


// CSRFValidationResponse represents a csrfvalidation_response
type CSRFValidationResponse struct {
	Valid bool `json:"valid,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
