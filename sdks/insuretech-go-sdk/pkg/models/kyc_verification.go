package models

import (
	"time"
)

// KYCVerification represents a kyc_verification
type KYCVerification struct {
	Documents string `json:"documents,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	Id string `json:"id"`
	EntityId string `json:"entity_id"`
	Method *VerificationMethod `json:"method"`
	Provider string `json:"provider,omitempty"`
	Status interface{} `json:"status"`
	VerificationResult string `json:"verification_result,omitempty"`
	Type *VerificationType `json:"type"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	EntityType string `json:"entity_type"`
	ProviderReference string `json:"provider_reference,omitempty"`
}
