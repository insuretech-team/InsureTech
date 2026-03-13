package models

import (
	"time"
)

// NIDInfo represents a nid_info
type NIDInfo struct {
	NidNumber string `json:"nid_number,omitempty"`
	FullName string `json:"full_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Verified bool `json:"verified,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	VerificationMethod string `json:"verification_method,omitempty"`
}
