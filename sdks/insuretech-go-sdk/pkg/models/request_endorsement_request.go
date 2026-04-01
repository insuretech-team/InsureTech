package models


// RequestEndorsementRequest represents a request_endorsement_request
type RequestEndorsementRequest struct {
	Changes string `json:"changes,omitempty"`
	EffectiveDate string `json:"effective_date,omitempty"`
	PolicyId string `json:"policy_id"`
	Reason string `json:"reason,omitempty"`
	Type string `json:"type"`
}
