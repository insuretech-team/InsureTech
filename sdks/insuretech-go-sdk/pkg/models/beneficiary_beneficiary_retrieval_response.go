package models


// BeneficiaryBeneficiaryRetrievalResponse represents a beneficiary_beneficiary_retrieval_response
type BeneficiaryBeneficiaryRetrievalResponse struct {
	BusinessDetails *BusinessBeneficiary `json:"business_details,omitempty"`
	Error *Error `json:"error,omitempty"`
	Beneficiary *Beneficiary `json:"beneficiary,omitempty"`
	IndividualDetails *IndividualBeneficiary `json:"individual_details,omitempty"`
}
