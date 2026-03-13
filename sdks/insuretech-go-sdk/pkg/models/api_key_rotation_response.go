package models

import (
	"time"
)

// APIKeyRotationResponse represents a api_key_rotation_response
type APIKeyRotationResponse struct {
	Error *Error `json:"error,omitempty"`
	NewKeyId string `json:"new_key_id,omitempty"`
	RawKey string `json:"raw_key,omitempty"`
	OldKeyId string `json:"old_key_id,omitempty"`
	OldKeyExpiresAt time.Time `json:"old_key_expires_at,omitempty"`
	Message string `json:"message,omitempty"`
}
