package models

import (
	"time"
)

// Applicant represents a applicant
type Applicant struct {
	AnnualIncome *Money `json:"annual_income,omitempty"`
	Address string `json:"address,omitempty"`
	HealthDeclaration *PolicyHealthDeclaration `json:"health_declaration,omitempty"`
	FullName string `json:"full_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	Occupation string `json:"occupation,omitempty"`
}
