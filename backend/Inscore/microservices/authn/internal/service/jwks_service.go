package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

// JWK represents a JSON Web Key (RFC 7517).
// NOTE: GetJWKS RPC was removed from auth_service.proto (API path conflict).
// JWKS is now served directly by the API gateway. JWKSService remains as an
// internal helper used by the gateway and integration tests.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSResult holds the set of public keys returned by GetJWKSInternal.
type JWKSResult struct {
	Keys []*JWK
}

// JWKSService exposes the public RSA key(s) in JWK Set format.
//
// NOTE: TokenService already implements GetJWKSInternal directly and AuthService
// delegates to it. JWKSService is a standalone helper for callers that need
// JWKS generation without a full TokenService (e.g. API gateways, integration
// tests, or CLI tooling that only holds the public key).
type JWKSService struct {
	publicKey *rsa.PublicKey
	keyID     string // kid claim in JWT header
}

// NewJWKSService creates a new JWKSService with the given RSA public key and
// key ID. publicKey may be nil; GetJWKSInternal will return an empty key set in that
// case.
func NewJWKSService(publicKey *rsa.PublicKey, keyID string) *JWKSService {
	return &JWKSService{
		publicKey: publicKey,
		keyID:     keyID,
	}
}

// GetJWKSInternal returns the JWK Set containing the RSA public key in RFC 7517 format.
//
// JWK structure:
//
//	{
//	  "kty": "RSA",
//	  "use": "sig",
//	  "alg": "RS256",
//	  "kid": "<keyID>",
//	  "n":   "<base64url(modulus)>",
//	  "e":   "<base64url(exponent)>"
//	}
func (s *JWKSService) GetJWKSInternal(_ context.Context) (*JWKSResult, error) {
	if s == nil || s.publicKey == nil {
		return &JWKSResult{Keys: []*JWK{}}, nil
	}

	// Modulus: base64url-encode the big-endian bytes (no padding).
	nBytes := s.publicKey.N.Bytes()
	nEncoded := base64.RawURLEncoding.EncodeToString(nBytes)

	// Exponent: convert to big-endian bytes, then base64url-encode.
	eBig := new(big.Int).SetInt64(int64(s.publicKey.E))
	eEncoded := base64.RawURLEncoding.EncodeToString(eBig.Bytes())

	jwk := &JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: s.keyID,
		N:   nEncoded,
		E:   eEncoded,
	}

	return &JWKSResult{
		Keys: []*JWK{jwk},
	}, nil
}
