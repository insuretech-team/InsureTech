package models


// DocumentVerification represents a document_verification
type DocumentVerification struct {
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	ExtractedData string `json:"extracted_data,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	KycVerificationId string `json:"kyc_verification_id"`
	DocumentType *KycDocumentType `json:"document_type"`
	DocumentNumber string `json:"document_number"`
	DocumentUrl string `json:"document_url,omitempty"`
	Status interface{} `json:"status"`
}
