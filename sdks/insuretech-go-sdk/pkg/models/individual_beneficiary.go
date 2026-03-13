package models

import (
	"time"
)

// IndividualBeneficiary represents a individual_beneficiary
type IndividualBeneficiary struct {
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`
	NomineeName string `json:"nominee_name,omitempty"`
	PermanentAddress *Address `json:"permanent_address,omitempty"`
	PresentAddress *Address `json:"present_address,omitempty"`
	FullName string `json:"full_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	BirthCertificateNumber string `json:"birth_certificate_number,omitempty"`
	AuditInfo *AuditInfo `json:"audit_info,omitempty"`
	MaritalStatus *MaritalStatus `json:"marital_status,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	FullNameBn string `json:"full_name_bn,omitempty"`
	Gender *BeneficiaryGender `json:"gender,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	NomineeRelationship string `json:"nominee_relationship,omitempty"`
	Occupation string `json:"occupation,omitempty"`
}
