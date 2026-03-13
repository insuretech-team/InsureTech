package models


// ClaimSubmissionRequest represents a claim_submission_request
type ClaimSubmissionRequest struct {
	IncidentDate string `json:"incident_date,omitempty"`
	IncidentDescription string `json:"incident_description,omitempty"`
	DocumentUrls []string `json:"document_urls,omitempty"`
	PolicyId string `json:"policy_id"`
	CustomerId string `json:"customer_id"`
	Type *ClaimType `json:"type"`
	ClaimedAmount *Money `json:"claimed_amount,omitempty"`
}
