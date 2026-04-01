package models


// MediaValidationResponse represents a media_validation_response
type MediaValidationResponse struct {
	ValidationErrors []string `json:"validation_errors,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
}
