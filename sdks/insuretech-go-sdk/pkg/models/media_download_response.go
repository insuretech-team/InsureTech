package models


// MediaDownloadResponse represents a media_download_response
type MediaDownloadResponse struct {
	DownloadUrl string `json:"download_url,omitempty"`
	ExpiresInSeconds string `json:"expires_in_seconds,omitempty"`
	Error *Error `json:"error,omitempty"`
}
