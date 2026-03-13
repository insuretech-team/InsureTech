package models

import (
	"time"
)

// PartnerVerificationResponse represents a partner_verification_response
type PartnerVerificationResponse struct {
	Verified bool `json:"verified,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
	Error *Error `json:"error,omitempty"`
}
