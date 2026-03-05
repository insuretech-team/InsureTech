package turn

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

// CredentialGenerator generates time-limited TURN credentials using HMAC
type CredentialGenerator struct {
	secret string
	ttl    time.Duration
}

// Credentials represents a set of TURN credentials
type Credentials struct {
	Username   string
	Credential string
	URLs       []string
	ExpiresAt  time.Time
}

// NewCredentialGenerator creates a new TURN credential generator
func NewCredentialGenerator(secret string, ttl time.Duration) *CredentialGenerator {
	return &CredentialGenerator{
		secret: secret,
		ttl:    ttl,
	}
}

// Generate creates time-limited TURN credentials for a user
// The credentials are valid for the configured TTL
// Format follows Coturn's time-limited credential mechanism
func (g *CredentialGenerator) Generate(username string) Credentials {
	expiry := time.Now().Add(g.ttl)
	timestamp := expiry.Unix()
	
	// Format: timestamp:username
	turnUsername := fmt.Sprintf("%d:%s", timestamp, username)
	
	// Generate HMAC-SHA1 credential
	h := hmac.New(sha1.New, []byte(g.secret))
	h.Write([]byte(turnUsername))
	credential := base64.StdEncoding.EncodeToString(h.Sum(nil))
	
	return Credentials{
		Username:   turnUsername,
		Credential: credential,
		ExpiresAt:  expiry,
	}
}

// GenerateWithURLs generates credentials with specific TURN server URLs
func (g *CredentialGenerator) GenerateWithURLs(username string, urls []string) Credentials {
	creds := g.Generate(username)
	creds.URLs = urls
	return creds
}

// IsValid checks if the credentials are still valid
func (c *Credentials) IsValid() bool {
	return time.Now().Before(c.ExpiresAt)
}

// ToICEServer converts credentials to ICE server configuration format
func (c *Credentials) ToICEServer() map[string]interface{} {
	server := map[string]interface{}{
		"urls":           c.URLs,
		"username":       c.Username,
		"credential":     c.Credential,
		"credentialType": "password",
	}
	return server
}
