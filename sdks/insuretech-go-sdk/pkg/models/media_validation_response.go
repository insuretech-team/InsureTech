package models


// MediaValidationResponse represents a media_validation_response
type MediaValidationResponse struct {
	Error *Error `json:"error,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
	ValidationErrors []string `json:"validation_errors,omitempty"`
}
