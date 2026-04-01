package models

import (
	"time"
)

// KYCVerification represents a kyc_verification
type KYCVerification struct {
	AuditInfo interface{} `json:"audit_info"`
	Documents string `json:"documents,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Id string `json:"id"`
	Method *VerificationMethod `json:"method"`
	Provider string `json:"provider,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	Status interface{} `json:"status"`
	Type *VerificationType `json:"type"`
	VerificationResult string `json:"verification_result,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
