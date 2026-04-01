package models


// DocumentDownloadResponse represents a document_download_response
type DocumentDownloadResponse struct {
	Content string `json:"content,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Filename string `json:"filename,omitempty"`
}
