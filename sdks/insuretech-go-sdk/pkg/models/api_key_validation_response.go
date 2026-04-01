package models


// ApiKeyValidationResponse represents a api_key_validation_response
type ApiKeyValidationResponse struct {
	ApiKeyId string `json:"api_key_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	OwnerType string `json:"owner_type,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	Valid bool `json:"valid,omitempty"`
}
