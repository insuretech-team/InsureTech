package models


// BeneficiaryIndividualBeneficiaryCreationRequest represents a beneficiary_individual_beneficiary_creation_request
type BeneficiaryIndividualBeneficiaryCreationRequest struct {
	UserId string `json:"user_id"`
	FullName string `json:"full_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Gender string `json:"gender,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email"`
	PartnerId string `json:"partner_id"`
}
