package models


// ClaimSubmissionRequest represents a claim_submission_request
type ClaimSubmissionRequest struct {
	ClaimedAmount *Money `json:"claimed_amount,omitempty"`
	CustomerId string `json:"customer_id"`
	DocumentUrls []string `json:"document_urls,omitempty"`
	IncidentDate string `json:"incident_date,omitempty"`
	IncidentDescription string `json:"incident_description,omitempty"`
	PolicyId string `json:"policy_id"`
	Type *ClaimType `json:"type"`
}
