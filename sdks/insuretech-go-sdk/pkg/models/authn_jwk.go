package models


// AuthnJWK represents a authn_jwk
type AuthnJWK struct {
	E string `json:"e,omitempty"`
	Kty string `json:"kty,omitempty"`
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N string `json:"n,omitempty"`
}
