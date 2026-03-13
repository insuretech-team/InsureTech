package models


// RevokeAPIKeyRequest represents a revoke_api_key_request
type RevokeAPIKeyRequest struct {
	KeyId string `json:"key_id"`
	Reason string `json:"reason,omitempty"`
}
