package models


// FileUpload represents a file_upload
type FileUpload struct {
	FileType *FileType `json:"file_type,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	IsPublic bool `json:"is_public,omitempty"`
	Content string `json:"content,omitempty"`
	Filename string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}
