package models


// AuthzJWK represents a authz_jwk
type AuthzJWK struct {
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`
	Kty string `json:"kty,omitempty"`
}
