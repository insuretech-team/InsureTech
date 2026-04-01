package models

import (
	"time"
)

// DocumentGeneration represents a document_generation
type DocumentGeneration struct {
	AuditInfo interface{} `json:"audit_info"`
	Data string `json:"data"`
	DocumentTemplateId string `json:"document_template_id"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	FileSizeBytes string `json:"file_size_bytes,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	GeneratedBy string `json:"generated_by,omitempty"`
	Id string `json:"id"`
	QrCodeData string `json:"qr_code_data,omitempty"`
	Status interface{} `json:"status"`
}
