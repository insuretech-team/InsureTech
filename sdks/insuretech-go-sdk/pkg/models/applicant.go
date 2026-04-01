package models

import (
	"time"
)

// Applicant represents a applicant
type Applicant struct {
	Address string `json:"address,omitempty"`
	AnnualIncome *Money `json:"annual_income,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	FullName string `json:"full_name,omitempty"`
	HealthDeclaration *PolicyHealthDeclaration `json:"health_declaration,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	Occupation string `json:"occupation,omitempty"`
}
