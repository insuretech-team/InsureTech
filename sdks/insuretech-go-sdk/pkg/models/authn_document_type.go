package models

import (
	"time"
)

// AuthnDocumentType represents a authn_document_type
type AuthnDocumentType struct {
	DocumentTypeId string `json:"document_type_id,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
