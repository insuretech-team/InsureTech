package models


// MediaUploadRequest represents a media_upload_request
type MediaUploadRequest struct {
	AutoOptimize bool `json:"auto_optimize,omitempty"`
	AutoThumbnail bool `json:"auto_thumbnail,omitempty"`
	FileId string `json:"file_id"`
	TenantId string `json:"tenant_id"`
	MediaType string `json:"media_type,omitempty"`
	EntityId string `json:"entity_id"`
	AutoValidate bool `json:"auto_validate,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	EntityType string `json:"entity_type"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}
