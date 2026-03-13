package models


// DownloadURLRetrievalRequest represents a download_url_retrieval_request
type DownloadURLRetrievalRequest struct {
	TenantId string `json:"tenant_id"`
	FileId string `json:"file_id"`
	ExpiresInMinutes int `json:"expires_in_minutes,omitempty"`
}
