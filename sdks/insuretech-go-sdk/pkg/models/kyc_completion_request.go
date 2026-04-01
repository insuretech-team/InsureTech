package models


// KYCCompletionRequest represents a kyc_completion_request
type KYCCompletionRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	NidBackUrl string `json:"nid_back_url,omitempty"`
	NidFrontUrl string `json:"nid_front_url,omitempty"`
	PorichoyVerificationId string `json:"porichoy_verification_id"`
	SelfieUrl string `json:"selfie_url,omitempty"`
}
