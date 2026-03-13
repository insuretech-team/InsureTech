package models


// DocumentDownloadResponse represents a document_download_response
type DocumentDownloadResponse struct {
	Filename string `json:"filename,omitempty"`
	Error *Error `json:"error,omitempty"`
	Content string `json:"content,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}
