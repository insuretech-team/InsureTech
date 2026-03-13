package models


// OptimizedDownloadResponse represents a optimized_download_response
type OptimizedDownloadResponse struct {
	DownloadUrl string `json:"download_url,omitempty"`
	ExpiresInSeconds string `json:"expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
