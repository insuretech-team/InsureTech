package models

import (
	"time"
)

// Nominee represents a nominee
type Nominee struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	FullName string `json:"full_name,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	NomineeDobText string `json:"nominee_dob_text,omitempty"`
	NomineeId string `json:"nominee_id,omitempty"`
	NomineeSharePercent float64 `json:"nominee_share_percent,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	SharePercentage float64 `json:"share_percentage,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
