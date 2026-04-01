package models


// AuthnJWK represents a authn_jwk
type AuthnJWK struct {
	Kty string `json:"kty,omitempty"`
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`
}
