package models


// DocumentVerification represents a document_verification
type DocumentVerification struct {
	AuditInfo interface{} `json:"audit_info"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	DocumentNumber string `json:"document_number"`
	DocumentType *KycDocumentType `json:"document_type"`
	DocumentUrl string `json:"document_url,omitempty"`
	ExtractedData string `json:"extracted_data,omitempty"`
	Id string `json:"id"`
	KycVerificationId string `json:"kyc_verification_id"`
	Status interface{} `json:"status"`
}
