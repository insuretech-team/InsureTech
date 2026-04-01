package models


// JWK represents a jwk
type JWK struct {
	Alg string `json:"alg,omitempty"`
	E string `json:"e,omitempty"`
	Kid string `json:"kid,omitempty"`
	Kty string `json:"kty,omitempty"`
	N string `json:"n,omitempty"`
	Use string `json:"use,omitempty"`
}
