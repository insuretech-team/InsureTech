package models

import (
	"time"
)

// AuthnDocumentType represents a authn_document_type
type AuthnDocumentType struct {
	Code string `json:"code,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Description string `json:"description,omitempty"`
	DocumentTypeId string `json:"document_type_id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
