package models


// DownloadURLRetrievalRequest represents a download_url_retrieval_request
type DownloadURLRetrievalRequest struct {
	ExpiresInMinutes int `json:"expires_in_minutes,omitempty"`
	FileId string `json:"file_id"`
	TenantId string `json:"tenant_id"`
}
