package models


// ThumbnailDownloadResponse represents a thumbnail_download_response
type ThumbnailDownloadResponse struct {
	DownloadUrl string `json:"download_url,omitempty"`
	ExpiresInSeconds string `json:"expires_in_seconds,omitempty"`
}
