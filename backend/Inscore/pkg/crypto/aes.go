package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const compactGCMNonceSize = 8

// EncryptPII encrypts a plaintext string using AES-GCM and returns a base64 encoded string.
// If the plaintext is empty, it returns the empty string.
func EncryptPII(plaintext string, key string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", errors.New("PII encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Partner PII is stored in legacy varchar-limited columns, so use a compact nonce size
	// for new writes while keeping decryption backward-compatible with older 12-byte nonce data.
	aesgcm, err := cipher.NewGCMWithNonceSize(block, compactGCMNonceSize)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPII decrypts a base64 encoded string using AES-GCM and returns the plaintext string.
// If the ciphertext is empty, it returns the empty string.
func DecryptPII(ciphertext string, key string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", errors.New("PII encryption key must be 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	nonceSizes := []int{12, compactGCMNonceSize}
	for _, nonceSize := range nonceSizes {
		aesgcm, err := cipher.NewGCMWithNonceSize(block, nonceSize)
		if err != nil {
			return "", err
		}
		if len(data) < nonceSize {
			continue
		}
		nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
		plaintext, err := aesgcm.Open(nil, nonce, ciphertextBytes, nil)
		if err == nil {
			return string(plaintext), nil
		}
	}

	return "", errors.New("failed to decrypt ciphertext")
}
