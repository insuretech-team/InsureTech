package models

import (
	"time"
)

// PartnerAPICredentialsRetrievalResponse represents a partner_api_credentials_retrieval_response
type PartnerAPICredentialsRetrievalResponse struct {
	ApiKey string `json:"api_key,omitempty"`
	ApiSecret string `json:"api_secret,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
