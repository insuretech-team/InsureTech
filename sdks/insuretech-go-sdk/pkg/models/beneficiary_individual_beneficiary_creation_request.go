package models


// BeneficiaryIndividualBeneficiaryCreationRequest represents a beneficiary_individual_beneficiary_creation_request
type BeneficiaryIndividualBeneficiaryCreationRequest struct {
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Email string `json:"email"`
	FullName string `json:"full_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	PartnerId string `json:"partner_id"`
	UserId string `json:"user_id"`
}
