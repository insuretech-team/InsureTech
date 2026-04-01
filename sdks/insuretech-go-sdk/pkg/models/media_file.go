package models


// MediaFile represents a media_file
type MediaFile struct {
	AuditInfo interface{} `json:"audit_info"`
	Dpi int `json:"dpi,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	FileId string `json:"file_id"`
	FileSizeBytes string `json:"file_size_bytes"`
	Height int `json:"height,omitempty"`
	Id string `json:"id"`
	MediaType *MediaType `json:"media_type"`
	MimeType string `json:"mime_type"`
	OcrText string `json:"ocr_text,omitempty"`
	OptimizedFileId string `json:"optimized_file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	ThumbnailFileId string `json:"thumbnail_file_id,omitempty"`
	UploadedBy string `json:"uploaded_by"`
	ValidationErrors string `json:"validation_errors,omitempty"`
	ValidationStatus interface{} `json:"validation_status"`
	VirusScanStatus interface{} `json:"virus_scan_status"`
	Width int `json:"width,omitempty"`
}
