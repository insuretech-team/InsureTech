package models

import (
	"time"
)

// Nominee represents a nominee
type Nominee struct {
	SharePercentage float64 `json:"share_percentage,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	NomineeDobText string `json:"nominee_dob_text,omitempty"`
	NomineeSharePercent float64 `json:"nominee_share_percent,omitempty"`
	NomineeId string `json:"nominee_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
