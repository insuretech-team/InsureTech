package models

import (
	"time"
)

// APIKeyCreationResponse represents a api_key_creation_response
type APIKeyCreationResponse struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	KeyId string `json:"key_id,omitempty"`
	Name string `json:"name,omitempty"`
	RawKey string `json:"raw_key,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}
