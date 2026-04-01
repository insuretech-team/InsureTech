package models


// InsurerInsurerCreationRequest represents a insurer_insurer_creation_request
type InsurerInsurerCreationRequest struct {
	Address string `json:"address,omitempty"`
	Code string `json:"code,omitempty"`
	Email string `json:"email"`
	Name string `json:"name"`
	NameBn string `json:"name_bn,omitempty"`
	Phone string `json:"phone,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	Type string `json:"type"`
}
