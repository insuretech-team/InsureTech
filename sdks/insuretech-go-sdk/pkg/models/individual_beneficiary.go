package models

import (
	"time"
)

// IndividualBeneficiary represents a individual_beneficiary
type IndividualBeneficiary struct {
	AuditInfo *AuditInfo `json:"audit_info,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	BirthCertificateNumber string `json:"birth_certificate_number,omitempty"`
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	FullName string `json:"full_name,omitempty"`
	FullNameBn string `json:"full_name_bn,omitempty"`
	Gender *BeneficiaryGender `json:"gender,omitempty"`
	MaritalStatus *MaritalStatus `json:"marital_status,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	NomineeName string `json:"nominee_name,omitempty"`
	NomineeRelationship string `json:"nominee_relationship,omitempty"`
	Occupation string `json:"occupation,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	PermanentAddress *Address `json:"permanent_address,omitempty"`
	PresentAddress *Address `json:"present_address,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
}
