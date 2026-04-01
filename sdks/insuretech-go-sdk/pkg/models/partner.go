package models

import (
	"time"
)

// Partner represents a partner
type Partner struct {
	AcquisitionCommissionRate float64 `json:"acquisition_commission_rate,omitempty"`
	BankAccount string `json:"bank_account,omitempty"`
	BankBranch string `json:"bank_branch,omitempty"`
	BankName string `json:"bank_name,omitempty"`
	Benefits *PartnerBenefits `json:"benefits,omitempty"`
	ClaimsAssistanceRate float64 `json:"claims_assistance_rate,omitempty"`
	Commission *CommissionStructure `json:"commission,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	FocalPersonId string `json:"focal_person_id,omitempty"`
	OnboardedAt time.Time `json:"onboarded_at,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	RenewalCommissionRate float64 `json:"renewal_commission_rate,omitempty"`
	Status *PartnerStatus `json:"status,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	TradeLicense string `json:"trade_license,omitempty"`
	Type *PartnerType `json:"type,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
