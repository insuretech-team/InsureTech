package models


// FileUploadResponse represents a file_upload_response
type FileUploadResponse struct {
	File *StoredFile `json:"file,omitempty"`
}
