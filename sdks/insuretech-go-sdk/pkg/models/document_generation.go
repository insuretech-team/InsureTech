package models

import (
	"time"
)

// DocumentGeneration represents a document_generation
type DocumentGeneration struct {
	EntityId string `json:"entity_id"`
	Status interface{} `json:"status"`
	GeneratedBy string `json:"generated_by,omitempty"`
	DocumentTemplateId string `json:"document_template_id"`
	EntityType string `json:"entity_type"`
	Data string `json:"data"`
	FileUrl string `json:"file_url,omitempty"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	QrCodeData string `json:"qr_code_data,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
}
