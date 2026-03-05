package service

import (
	"crypto/sha256"
	"encoding/hex"
)

// sessionTokenLookup computes a deterministic SHA-256 hex string used to look up server-side sessions.
// The plain token is still verified against bcrypt hash for security.
func sessionTokenLookup(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
