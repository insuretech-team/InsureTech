package models


// ClaimEvaluationRequest represents a claim_evaluation_request
type ClaimEvaluationRequest struct {
	ClaimId string `json:"claim_id"`
	PolicyId string `json:"policy_id"`
	DocumentUrls []string `json:"document_urls,omitempty"`
}
