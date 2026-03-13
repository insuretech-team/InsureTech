package models

import (
	"time"
)

// UserDocument represents a user_document
type UserDocument struct {
	UserId string `json:"user_id,omitempty"`
	DocumentTypeId string `json:"document_type_id,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	FileUrl string `json:"file_url,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UserDocumentId string `json:"user_document_id,omitempty"`
}
