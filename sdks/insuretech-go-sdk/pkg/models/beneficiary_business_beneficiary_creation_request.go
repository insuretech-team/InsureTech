package models


// BeneficiaryBusinessBeneficiaryCreationRequest represents a beneficiary_business_beneficiary_creation_request
type BeneficiaryBusinessBeneficiaryCreationRequest struct {
	BusinessName string `json:"business_name,omitempty"`
	FocalPersonMobile string `json:"focal_person_mobile,omitempty"`
	FocalPersonName string `json:"focal_person_name,omitempty"`
	PartnerId string `json:"partner_id"`
	TinNumber string `json:"tin_number,omitempty"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	UserId string `json:"user_id"`
}
