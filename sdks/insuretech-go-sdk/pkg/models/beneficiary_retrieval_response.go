package models


// BeneficiaryRetrievalResponse represents a beneficiary_retrieval_response
type BeneficiaryRetrievalResponse struct {
	IndividualDetails *IndividualBeneficiary `json:"individual_details,omitempty"`
	BusinessDetails *BusinessBeneficiary `json:"business_details,omitempty"`
	Error *Error `json:"error,omitempty"`
	Beneficiary *Beneficiary `json:"beneficiary,omitempty"`
}
