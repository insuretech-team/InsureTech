package models


// DocumentGenerationResponse represents a document_generation_response
type DocumentGenerationResponse struct {
	DocumentId string `json:"document_id,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
