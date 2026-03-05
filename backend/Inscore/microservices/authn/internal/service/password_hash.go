package service

// password_hash.go — Argon2id password hashing with bcrypt migration path.
// Decision: new passwords use argon2id. On login, if bcrypt hash detected → verify with bcrypt,
// then re-hash with argon2id for next login (transparent migration).
//
// Argon2id parameters (OWASP recommended):
//   memory=64MB, iterations=3, parallelism=4, saltLen=16, keyLen=32
//
// Hash format: $argon2id$v=19$m=65536,t=3,p=4$<salt_b64>$<hash_b64>

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	argon2Memory      = 64 * 1024 // 64 MB
	argon2Iterations  = 3
	argon2Parallelism = 4
	argon2SaltLen     = 16
	argon2KeyLen      = 32
)

// hashPassword hashes a password using Argon2id.
func hashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		logger.Errorf("generate salt: %v", err)
		return "", errors.New("generate salt")
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)
	encoded := "$argon2id$v=19$m=" + strconv.Itoa(argon2Memory) + ",t=" + strconv.Itoa(argon2Iterations) + ",p=" + strconv.Itoa(argon2Parallelism) + "$" +
		base64.RawStdEncoding.EncodeToString(salt) + "$" +
		base64.RawStdEncoding.EncodeToString(hash)
	return encoded, nil
}

// verifyPassword checks a password against a stored hash.
// Supports both argon2id (new) and bcrypt (legacy migration path).
// Returns (valid, needsRehash, error). needsRehash=true when bcrypt hash detected (caller should re-hash).
func verifyPassword(password, encodedHash string) (valid bool, needsRehash bool, err error) {
	if strings.HasPrefix(encodedHash, "$argon2id$") {
		ok, err := verifyArgon2id(password, encodedHash)
		return ok, false, err
	}
	// bcrypt (legacy)
	err = bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
	if err != nil {
		return false, false, nil
	}
	return true, true, nil // valid + needs re-hash to argon2id
}

func verifyArgon2id(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	// $argon2id$v=19$m=...,t=...,p=...$<salt>$<hash>
	if len(parts) != 6 {
		return false, errors.New("invalid argon2id hash format")
	}
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, errors.New("invalid argon2id params")
	}
	memStr := strings.TrimPrefix(params[0], "m=")
	tStr := strings.TrimPrefix(params[1], "t=")
	pStr := strings.TrimPrefix(params[2], "p=")
	mem, err1 := strconv.ParseUint(memStr, 10, 32)
	it, err2 := strconv.ParseUint(tStr, 10, 32)
	p, err3 := strconv.ParseUint(pStr, 10, 32)
	if err1 != nil || err2 != nil || err3 != nil {
		return false, errors.New("parse argon2id params failed")
	}
	memory := uint32(mem)
	iterations := uint32(it)
	parallelism := uint8(p)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		logger.Errorf("decode salt: %v", err)
		return false, errors.New("decode salt")
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		logger.Errorf("decode hash: %v", err)
		return false, errors.New("decode hash")
	}
	keyLen := uint32(len(storedHash))
	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)
	if subtle.ConstantTimeCompare(computedHash, storedHash) != 1 {
		return false, nil
	}
	return true, nil
}
