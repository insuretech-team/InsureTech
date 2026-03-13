package models

import (
	"time"
)

// TokenConfig represents a token_config
type TokenConfig struct {
	Kid string `json:"kid"`
	Algorithm string `json:"algorithm"`
	PublicKeyPem string `json:"public_key_pem"`
	PrivateKeyRef string `json:"private_key_ref"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	RotatedAt time.Time `json:"rotated_at,omitempty"`
}
