package models

import (
	"time"
)

// Partner represents a partner
type Partner struct {
	Benefits *PartnerBenefits `json:"benefits,omitempty"`
	ClaimsAssistanceRate float64 `json:"claims_assistance_rate,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	Status *PartnerStatus `json:"status,omitempty"`
	BankName string `json:"bank_name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	AcquisitionCommissionRate float64 `json:"acquisition_commission_rate,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Commission *CommissionStructure `json:"commission,omitempty"`
	Type *PartnerType `json:"type,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	RenewalCommissionRate float64 `json:"renewal_commission_rate,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	TradeLicense string `json:"trade_license,omitempty"`
	BankAccount string `json:"bank_account,omitempty"`
	BankBranch string `json:"bank_branch,omitempty"`
	OnboardedAt time.Time `json:"onboarded_at,omitempty"`
	FocalPersonId string `json:"focal_person_id,omitempty"`
}
