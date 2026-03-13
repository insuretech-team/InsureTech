package models


// TokenValidationRequest represents a token_validation_request
type TokenValidationRequest struct {
	AccessToken string `json:"access_token,omitempty"`
	SessionId string `json:"session_id"`
}
