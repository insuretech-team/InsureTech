package models


// DisputeClaimRequest represents a dispute_claim_request
type DisputeClaimRequest struct {
	ClaimId string `json:"claim_id"`
	CustomerId string `json:"customer_id"`
	DisputeReason string `json:"dispute_reason,omitempty"`
	SupportingDocumentUrls []string `json:"supporting_document_urls,omitempty"`
}
