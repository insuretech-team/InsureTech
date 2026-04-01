package models


// FileUpload represents a file_upload
type FileUpload struct {
	Content string `json:"content,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	FileType *FileType `json:"file_type,omitempty"`
	Filename string `json:"filename,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
}
