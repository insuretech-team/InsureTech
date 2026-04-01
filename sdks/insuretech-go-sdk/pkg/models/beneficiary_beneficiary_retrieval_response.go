package models


// BeneficiaryBeneficiaryRetrievalResponse represents a beneficiary_beneficiary_retrieval_response
type BeneficiaryBeneficiaryRetrievalResponse struct {
	Beneficiary *Beneficiary `json:"beneficiary,omitempty"`
	BusinessDetails *BusinessBeneficiary `json:"business_details,omitempty"`
	IndividualDetails *IndividualBeneficiary `json:"individual_details,omitempty"`
}
