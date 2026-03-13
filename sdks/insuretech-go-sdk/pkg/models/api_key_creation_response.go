package models

import (
	"time"
)

// APIKeyCreationResponse represents a api_key_creation_response
type APIKeyCreationResponse struct {
	RawKey string `json:"raw_key,omitempty"`
	Name string `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error *Error `json:"error,omitempty"`
	KeyId string `json:"key_id,omitempty"`
}
