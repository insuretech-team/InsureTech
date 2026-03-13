package models


// RequestEndorsementRequest represents a request_endorsement_request
type RequestEndorsementRequest struct {
	Reason string `json:"reason,omitempty"`
	Changes string `json:"changes,omitempty"`
	EffectiveDate string `json:"effective_date,omitempty"`
	PolicyId string `json:"policy_id"`
	Type string `json:"type"`
}
