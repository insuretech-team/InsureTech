package models


// JWK represents a jwk
type JWK struct {
	Kty string `json:"kty,omitempty"`
	Use string `json:"use,omitempty"`
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`
}
