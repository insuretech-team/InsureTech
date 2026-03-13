package models


// MediaFile represents a media_file
type MediaFile struct {
	Dpi int `json:"dpi,omitempty"`
	ValidationStatus interface{} `json:"validation_status"`
	Id string `json:"id"`
	TenantId string `json:"tenant_id,omitempty"`
	MediaType *MediaType `json:"media_type"`
	EntityId string `json:"entity_id,omitempty"`
	OcrText string `json:"ocr_text,omitempty"`
	MimeType string `json:"mime_type"`
	OptimizedFileId string `json:"optimized_file_id,omitempty"`
	ThumbnailFileId string `json:"thumbnail_file_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Height int `json:"height,omitempty"`
	ValidationErrors string `json:"validation_errors,omitempty"`
	VirusScanStatus interface{} `json:"virus_scan_status"`
	UploadedBy string `json:"uploaded_by"`
	AuditInfo interface{} `json:"audit_info"`
	FileId string `json:"file_id"`
	FileSizeBytes string `json:"file_size_bytes"`
	Width int `json:"width,omitempty"`
}
