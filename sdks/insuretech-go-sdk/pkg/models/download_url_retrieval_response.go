package models

import (
	"time"
)

// DownloadURLRetrievalResponse represents a download_url_retrieval_response
type DownloadURLRetrievalResponse struct {
	DownloadUrl string `json:"download_url,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
