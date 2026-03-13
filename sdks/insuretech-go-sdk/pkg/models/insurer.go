package models

import (
	"time"
)

// Insurer represents a insurer
type Insurer struct {
	HeadOfficeAddress interface{} `json:"head_office_address"`
	WebsiteUrl string `json:"website_url,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	LogoUrl string `json:"logo_url,omitempty"`
	Id string `json:"id"`
	Name string `json:"name"`
	NameBn string `json:"name_bn,omitempty"`
	Type *InsurerType `json:"type"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	IdraLicenseExpiry time.Time `json:"idra_license_expiry,omitempty"`
	ContactInfo interface{} `json:"contact_info"`
	Status interface{} `json:"status"`
	TinNumber string `json:"tin_number,omitempty"`
	FinancialRating string `json:"financial_rating,omitempty"`
	PaidUpCapital *Money `json:"paid_up_capital,omitempty"`
	RegisteredAddress interface{} `json:"registered_address"`
	Code string `json:"code"`
	IdraLicenseNumber string `json:"idra_license_number,omitempty"`
}
