package models


// RiskScoreUpdateRequest represents a risk_score_update_request
type RiskScoreUpdateRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Reason string `json:"reason,omitempty"`
	RiskScore string `json:"risk_score,omitempty"`
}
