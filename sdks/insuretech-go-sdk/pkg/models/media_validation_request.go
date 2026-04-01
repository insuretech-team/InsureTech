package models


// MediaValidationRequest represents a media_validation_request
type MediaValidationRequest struct {
	MediaId string `json:"media_id"`
	ValidationRules []string `json:"validation_rules,omitempty"`
}
