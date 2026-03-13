package models


// RiskScoreUpdateResponse represents a risk_score_update_response
type RiskScoreUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
