package models


// FileRetrievalResponse represents a file_retrieval_response
type FileRetrievalResponse struct {
	File *StoredFile `json:"file,omitempty"`
}
