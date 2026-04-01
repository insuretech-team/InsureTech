package models


// DocumentGenerationResponse represents a document_generation_response
type DocumentGenerationResponse struct {
	ContentType string `json:"content_type,omitempty"`
	DocumentId string `json:"document_id,omitempty"`
	FileBytes string `json:"file_bytes,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	Filename string `json:"filename,omitempty"`
}
