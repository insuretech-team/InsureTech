package models


// ApiKeyGenerationResponse represents a api_key_generation_response
type ApiKeyGenerationResponse struct {
	ApiKey string `json:"api_key,omitempty"`
	ApiKeyId string `json:"api_key_id,omitempty"`
}
