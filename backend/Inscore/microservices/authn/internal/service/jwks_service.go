package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"math/big"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
)

// JWKSService exposes the public RSA key(s) in JWK Set format.
//
// NOTE: TokenService already implements GetJWKS directly and AuthService
// delegates to it. JWKSService is a standalone helper for callers that need
// JWKS generation without a full TokenService (e.g. API gateways, integration
// tests, or CLI tooling that only holds the public key).
type JWKSService struct {
	publicKey *rsa.PublicKey
	keyID     string // kid claim in JWT header
}

// NewJWKSService creates a new JWKSService with the given RSA public key and
// key ID. publicKey may be nil; GetJWKS will return an empty key set in that
// case.
func NewJWKSService(publicKey *rsa.PublicKey, keyID string) *JWKSService {
	return &JWKSService{
		publicKey: publicKey,
		keyID:     keyID,
	}
}

// GetJWKS returns the JWK Set containing the RSA public key in RFC 7517 format.
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
func (s *JWKSService) GetJWKS(ctx context.Context, req *authnservicev1.GetJWKSRequest) (*authnservicev1.GetJWKSResponse, error) {
	if s == nil || s.publicKey == nil {
		return &authnservicev1.GetJWKSResponse{Keys: []*authnservicev1.JWK{}}, nil
	}

	// Modulus: base64url-encode the big-endian bytes (no padding).
	nBytes := s.publicKey.N.Bytes()
	nEncoded := base64.RawURLEncoding.EncodeToString(nBytes)

	// Exponent: convert to big-endian bytes, then base64url-encode.
	eBig := new(big.Int).SetInt64(int64(s.publicKey.E))
	eEncoded := base64.RawURLEncoding.EncodeToString(eBig.Bytes())

	jwk := &authnservicev1.JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: s.keyID,
		N:   nEncoded,
		E:   eEncoded,
	}

	return &authnservicev1.GetJWKSResponse{
		Keys: []*authnservicev1.JWK{jwk},
	}, nil
}
