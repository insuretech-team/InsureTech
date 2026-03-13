package models


// FilesUploadResponse represents a files_upload_response
type FilesUploadResponse struct {
	Files []*StoredFile `json:"files,omitempty"`
}
