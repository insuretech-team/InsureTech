package service

// argon2id.go - Exported Argon2id helpers.
//
// The core hashing logic lives in password_hash.go (unexported hashPassword /
// verifyPassword). This file exposes a public API with explicit parameter
// control and adds IsArgon2idHash for migration detection.

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"golang.org/x/crypto/argon2"
)

// Argon2idParams holds the Argon2id configuration.
// The defaults match OWASP's 2024 minimum recommended values.
type Argon2idParams struct {
	Memory      uint32 // kilobytes of RAM (default 64*1024 = 64 MB)
	Iterations  uint32 // number of passes (default 3)
	Parallelism uint8  // degree of parallelism (default 4)
	SaltLength  uint32 // random salt length in bytes (default 16)
	KeyLength   uint32 // derived key length in bytes (default 32)
}

// DefaultArgon2idParams are the recommended production parameters.
var DefaultArgon2idParams = &Argon2idParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword hashes a plaintext password using Argon2id and returns the
// PHC-style encoded string:
//
//	$argon2id$v=19$m=<mem>,t=<iter>,p=<par>$<salt_b64>$<hash_b64>
//
// Pass nil for p to use DefaultArgon2idParams.
func HashPassword(password string, p *Argon2idParams) (string, error) {
	if p == nil {
		p = DefaultArgon2idParams
	}
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		logger.Errorf("argon2id: generate salt: %v", err)
		return "", errors.New("argon2id: generate salt")
	}
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks a plaintext password against an Argon2id-encoded hash.
// Returns true when the password matches.
//
// Only Argon2id hashes (strings starting with "$argon2id$") are accepted; for
// bcrypt migration use the lower-level verifyPassword in password_hash.go.
func VerifyPassword(password, encodedHash string) (bool, error) {
	if !IsArgon2idHash(encodedHash) {
		logger.Errorf("argon2id: unsupported hash format (use verifyPassword for bcrypt migration)")
		return false, errors.New("argon2id: unsupported hash format (use verifyPassword for bcrypt migration)")
	}

	parts := strings.Split(encodedHash, "$")
	// Expected: ["", "argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 6 {
		return false, fmt.Errorf("argon2id: invalid hash format (expected 6 parts, got %d)", len(parts))
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		logger.Errorf("argon2id: parse params: %v", err)
		return false, errors.New("argon2id: parse params")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		logger.Errorf("argon2id: decode salt: %v", err)
		return false, errors.New("argon2id: decode salt")
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		logger.Errorf("argon2id: decode hash: %v", err)
		return false, errors.New("argon2id: decode hash")
	}

	keyLen := uint32(len(storedHash))
	computed := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)
	if subtle.ConstantTimeCompare(computed, storedHash) != 1 {
		return false, nil
	}
	return true, nil
}

// IsArgon2idHash reports whether the encoded string is an Argon2id hash
// produced by HashPassword or hashPassword. Useful for bcrypt migration
// detection at the call site.
func IsArgon2idHash(encoded string) bool {
	return strings.HasPrefix(encoded, "$argon2id$")
}
