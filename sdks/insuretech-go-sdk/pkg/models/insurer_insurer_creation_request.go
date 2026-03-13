package models


// InsurerInsurerCreationRequest represents a insurer_insurer_creation_request
type InsurerInsurerCreationRequest struct {
	Code string `json:"code,omitempty"`
	NameBn string `json:"name_bn,omitempty"`
	Type string `json:"type"`
	Email string `json:"email"`
	Address string `json:"address,omitempty"`
	Name string `json:"name"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	Phone string `json:"phone,omitempty"`
}
