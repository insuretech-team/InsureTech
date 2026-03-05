// Package pii provides field-level AES-256-GCM encryption and HMAC-SHA256
// blind-index helpers for PII data at rest.
//
// Usage pattern:
//
//	enc := pii.NewEncryptor(aesKey, hmacKey)
//	ciphertext, err := enc.Encrypt("01711000000")
//	idx          := enc.BlindIndex("01711000000")   // deterministic; safe for WHERE-clause lookup
//	plaintext, err := enc.Decrypt(ciphertext)
//
// Key sizes:
//   - AES key:  exactly 32 bytes (AES-256)
//   - HMAC key: exactly 32 bytes (HMAC-SHA256)
//
// Both keys should be loaded from environment variables or a secret manager
// and MUST NOT be stored in code or version control.
package pii

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
)

// ErrEmptyPlaintext is returned when encrypting an empty string.
var ErrEmptyPlaintext = errors.New("pii: plaintext must not be empty")

// Encryptor holds the AES-256-GCM and HMAC-SHA256 keys.
type Encryptor struct {
	aesKey  []byte // 32 bytes
	hmacKey []byte // 32 bytes
}

// NewEncryptor creates an Encryptor with the provided keys.
// Both keys must be exactly 32 bytes; an error is returned otherwise.
func NewEncryptor(aesKey, hmacKey []byte) (*Encryptor, error) {
	if len(aesKey) != 32 {
		return nil, fmt.Errorf("pii: AES key must be 32 bytes, got %d", len(aesKey))
	}
	if len(hmacKey) != 32 {
		return nil, fmt.Errorf("pii: HMAC key must be 32 bytes, got %d", len(hmacKey))
	}
	return &Encryptor{
		aesKey:  aesKey,
		hmacKey: hmacKey,
	}, nil
}

// NewEncryptorFromEnv creates an Encryptor from environment variables.
// Expected env vars:
//   - PII_AES_KEY:  64 hex chars (32 bytes)
//   - PII_HMAC_KEY: 64 hex chars (32 bytes)
func NewEncryptorFromEnv() (*Encryptor, error) {
	aesHex := os.Getenv("PII_AES_KEY")
	hmacHex := os.Getenv("PII_HMAC_KEY")

	if aesHex == "" {
		logger.Errorf("pii: PII_AES_KEY environment variable not set")
		return nil, errors.New("pii: PII_AES_KEY environment variable not set")
	}
	if hmacHex == "" {
		logger.Errorf("pii: PII_HMAC_KEY environment variable not set")
		return nil, errors.New("pii: PII_HMAC_KEY environment variable not set")
	}

	aesKey, err := hex.DecodeString(aesHex)
	if err != nil {
		logger.Errorf("pii: PII_AES_KEY is not valid hex: %v", err)
		return nil, errors.New("pii: PII_AES_KEY is not valid hex")
	}
	hmacKey, err := hex.DecodeString(hmacHex)
	if err != nil {
		logger.Errorf("pii: PII_HMAC_KEY is not valid hex: %v", err)
		return nil, errors.New("pii: PII_HMAC_KEY is not valid hex")
	}

	return NewEncryptor(aesKey, hmacKey)
}

// Encrypt encrypts plaintext with AES-256-GCM and returns a base64url-encoded
// string in the format: base64(nonce || ciphertext || tag).
// A fresh random 12-byte nonce is generated for every call (nonce reuse would
// be catastrophic for GCM; this design avoids it entirely).
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", ErrEmptyPlaintext
	}

	block, err := aes.NewCipher(e.aesKey)
	if err != nil {
		logger.Errorf("pii: failed to create AES cipher: %v", err)
		return "", errors.New("pii: failed to create AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Errorf("pii: failed to create GCM: %v", err)
		return "", errors.New("pii: failed to create GCM")
	}

	nonce := make([]byte, gcm.NonceSize()) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Errorf("pii: failed to generate nonce: %v", err)
		return "", errors.New("pii: failed to generate nonce")
	}

	// Seal appends ciphertext+tag after nonce.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64url-encoded ciphertext produced by Encrypt.
func (e *Encryptor) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Errorf("pii: failed to base64-decode ciphertext: %v", err)
		return "", errors.New("pii: failed to base64-decode ciphertext")
	}

	block, err := aes.NewCipher(e.aesKey)
	if err != nil {
		logger.Errorf("pii: failed to create AES cipher: %v", err)
		return "", errors.New("pii: failed to create AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Errorf("pii: failed to create GCM: %v", err)
		return "", errors.New("pii: failed to create GCM")
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		logger.Errorf("pii: ciphertext too short")
		return "", errors.New("pii: ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Errorf("pii: decryption failed (wrong key or corrupt data): %v", err)
		return "", errors.New("pii: decryption failed (wrong key or corrupt data)")
	}

	return string(plaintext), nil
}

// BlindIndex returns a deterministic HMAC-SHA256 hex digest of the plaintext,
// suitable for use as a database index column to enable equality lookups
// without exposing the plaintext. The result is always 64 hex characters.
//
// This MUST use a separate key from the AES key to prevent related-key attacks.
func (e *Encryptor) BlindIndex(plaintext string) string {
	if plaintext == "" {
		return ""
	}
	mac := hmac.New(sha256.New, e.hmacKey)
	mac.Write([]byte(plaintext))
	return hex.EncodeToString(mac.Sum(nil))
}

// EncryptIfNonEmpty encrypts plaintext only if it is non-empty; returns "" for
// empty input without error. Useful for optional fields like NID.
func (e *Encryptor) EncryptIfNonEmpty(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	return e.Encrypt(plaintext)
}

// DecryptIfNonEmpty decrypts encoded only if it is non-empty; returns "" for
// empty input without error. Mirrors EncryptIfNonEmpty.
func (e *Encryptor) DecryptIfNonEmpty(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	return e.Decrypt(encoded)
}

// GenerateBiometricToken generates a cryptographically secure 256-bit (32 byte)
// random token, returns:
//   - raw:       the plaintext token (send once to device; never store plaintext)
//   - encrypted: AES-256-GCM ciphertext of the token (store in DB)
//   - lookupIdx: HMAC-SHA256 blind index (store in DB for lookup)
func (e *Encryptor) GenerateBiometricToken() (raw, encrypted, lookupIdx string, err error) {
	b := make([]byte, 32) // 256-bit token
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		logger.Errorf("pii: failed to generate biometric token random bytes: %v", err)
		err = errors.New("pii: failed to generate biometric token random bytes")
		return
	}
	raw = hex.EncodeToString(b)

	encrypted, err = e.Encrypt(raw)
	if err != nil {
		logger.Errorf("pii: failed to encrypt biometric token: %v", err)
		err = errors.New("pii: failed to encrypt biometric token")
		return
	}

	lookupIdx = e.BlindIndex(raw)
	return
}
