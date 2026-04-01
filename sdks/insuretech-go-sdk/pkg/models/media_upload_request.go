package models


// MediaUploadRequest represents a media_upload_request
type MediaUploadRequest struct {
	AutoOptimize bool `json:"auto_optimize,omitempty"`
	AutoThumbnail bool `json:"auto_thumbnail,omitempty"`
	AutoValidate bool `json:"auto_validate,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	FileId string `json:"file_id"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	MediaType string `json:"media_type,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	TenantId string `json:"tenant_id"`
	UploadedBy string `json:"uploaded_by,omitempty"`
}
