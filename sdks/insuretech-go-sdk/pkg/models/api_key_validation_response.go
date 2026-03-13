package models


// ApiKeyValidationResponse represents a api_key_validation_response
type ApiKeyValidationResponse struct {
	OwnerType string `json:"owner_type,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Valid bool `json:"valid,omitempty"`
	ApiKeyId string `json:"api_key_id,omitempty"`
}
