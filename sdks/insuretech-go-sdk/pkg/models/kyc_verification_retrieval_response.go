package models


// KYCVerificationRetrievalResponse represents a kyc_verification_retrieval_response
type KYCVerificationRetrievalResponse struct {
	Documents []*DocumentVerification `json:"documents,omitempty"`
	KycVerification *KYCVerification `json:"kyc_verification,omitempty"`
}
