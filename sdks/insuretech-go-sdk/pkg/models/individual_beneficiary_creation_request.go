package models


// IndividualBeneficiaryCreationRequest represents a individual_beneficiary_creation_request
type IndividualBeneficiaryCreationRequest struct {
	Gender string `json:"gender,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email"`
	PartnerId string `json:"partner_id"`
	UserId string `json:"user_id"`
	FullName string `json:"full_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}
