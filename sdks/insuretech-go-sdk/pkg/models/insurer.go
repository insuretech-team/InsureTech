package models

import (
	"time"
)

// Insurer represents a insurer
type Insurer struct {
	AuditInfo interface{} `json:"audit_info"`
	Code string `json:"code"`
	ContactInfo interface{} `json:"contact_info"`
	FinancialRating string `json:"financial_rating,omitempty"`
	HeadOfficeAddress interface{} `json:"head_office_address"`
	Id string `json:"id"`
	IdraLicenseExpiry time.Time `json:"idra_license_expiry,omitempty"`
	IdraLicenseNumber string `json:"idra_license_number,omitempty"`
	LogoUrl string `json:"logo_url,omitempty"`
	Name string `json:"name"`
	NameBn string `json:"name_bn,omitempty"`
	PaidUpCapital *Money `json:"paid_up_capital,omitempty"`
	RegisteredAddress interface{} `json:"registered_address"`
	Status interface{} `json:"status"`
	TinNumber string `json:"tin_number,omitempty"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	Type *InsurerType `json:"type"`
	WebsiteUrl string `json:"website_url,omitempty"`
}
