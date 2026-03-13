package models


// RevokeAPIKeyResponse represents a revoke_api_key_response
type RevokeAPIKeyResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
