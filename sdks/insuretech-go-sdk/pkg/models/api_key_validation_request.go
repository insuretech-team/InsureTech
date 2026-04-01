package models


// ApiKeyValidationRequest represents a api_key_validation_request
type ApiKeyValidationRequest struct {
	ApiKey string `json:"api_key"`
	RequestIp string `json:"request_ip,omitempty"`
	RequiredScope string `json:"required_scope,omitempty"`
}
