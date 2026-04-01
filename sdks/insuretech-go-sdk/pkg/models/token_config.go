package models

import (
	"time"
)

// TokenConfig represents a token_config
type TokenConfig struct {
	Algorithm string `json:"algorithm"`
	CreatedAt time.Time `json:"created_at"`
	IsActive bool `json:"is_active"`
	Kid string `json:"kid"`
	PrivateKeyRef string `json:"private_key_ref"`
	PublicKeyPem string `json:"public_key_pem"`
	RotatedAt time.Time `json:"rotated_at,omitempty"`
}
