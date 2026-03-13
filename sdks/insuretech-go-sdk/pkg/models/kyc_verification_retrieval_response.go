package models


// KYCVerificationRetrievalResponse represents a kyc_verification_retrieval_response
type KYCVerificationRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	KycVerification *KYCVerification `json:"kyc_verification,omitempty"`
	Documents []*DocumentVerification `json:"documents,omitempty"`
}
