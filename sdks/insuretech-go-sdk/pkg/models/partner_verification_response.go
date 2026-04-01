package models

import (
	"time"
)

// PartnerVerificationResponse represents a partner_verification_response
type PartnerVerificationResponse struct {
	VerificationStatus string `json:"verification_status,omitempty"`
	Verified bool `json:"verified,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
}
