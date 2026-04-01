package models

import (
	"time"
)

// BusinessBeneficiary represents a business_beneficiary
type BusinessBeneficiary struct {
	ActivePoliciesCount int `json:"active_policies_count,omitempty"`
	AuditInfo *AuditInfo `json:"audit_info,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	BinNumber string `json:"bin_number,omitempty"`
	BusinessAddress *Address `json:"business_address,omitempty"`
	BusinessName string `json:"business_name,omitempty"`
	BusinessNameBn string `json:"business_name_bn,omitempty"`
	BusinessType *BusinessType `json:"business_type,omitempty"`
	ContactInfo *ContactInfo `json:"contact_info,omitempty"`
	EmployeeCount int `json:"employee_count,omitempty"`
	FocalPersonContact *ContactInfo `json:"focal_person_contact,omitempty"`
	FocalPersonDesignation string `json:"focal_person_designation,omitempty"`
	FocalPersonName string `json:"focal_person_name,omitempty"`
	FocalPersonNid string `json:"focal_person_nid,omitempty"`
	Id string `json:"id,omitempty"`
	IncorporationDate time.Time `json:"incorporation_date,omitempty"`
	IndustrySector string `json:"industry_sector,omitempty"`
	PendingActionsCount int `json:"pending_actions_count,omitempty"`
	PrimaryContact *PrimaryContact `json:"primary_contact,omitempty"`
	RegisteredAddress *Address `json:"registered_address,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	TaxId string `json:"tax_id,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	TotalEmployeesCovered int `json:"total_employees_covered,omitempty"`
	TotalPremiumAmount *Money `json:"total_premium_amount,omitempty"`
	TradeLicenseExpiryDate time.Time `json:"trade_license_expiry_date,omitempty"`
	TradeLicenseIssueDate time.Time `json:"trade_license_issue_date,omitempty"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
}
