package models


// InsurerCreationRequest represents a insurer_creation_request
type InsurerCreationRequest struct {
	NameBn string `json:"name_bn,omitempty"`
	Type string `json:"type"`
	Email string `json:"email"`
	Code string `json:"code,omitempty"`
	Name string `json:"name"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	Phone string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}
