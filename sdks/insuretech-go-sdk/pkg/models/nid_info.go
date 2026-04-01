package models

import (
	"time"
)

// NIDInfo represents a nid_info
type NIDInfo struct {
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	FullName string `json:"full_name,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	VerificationMethod string `json:"verification_method,omitempty"`
	Verified bool `json:"verified,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
}
