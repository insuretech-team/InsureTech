package models


// BeneficiaryBusinessBeneficiaryCreationRequest represents a beneficiary_business_beneficiary_creation_request
type BeneficiaryBusinessBeneficiaryCreationRequest struct {
	UserId string `json:"user_id"`
	BusinessName string `json:"business_name,omitempty"`
	TradeLicenseNumber string `json:"trade_license_number,omitempty"`
	TinNumber string `json:"tin_number,omitempty"`
	FocalPersonName string `json:"focal_person_name,omitempty"`
	FocalPersonMobile string `json:"focal_person_mobile,omitempty"`
	PartnerId string `json:"partner_id"`
}
