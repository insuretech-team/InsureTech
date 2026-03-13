package models


// MediaValidationRequest represents a media_validation_request
type MediaValidationRequest struct {
	ValidationRules []string `json:"validation_rules,omitempty"`
	MediaId string `json:"media_id"`
}
