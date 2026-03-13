package models


// DetectFraudRequest represents a detect_fraud_request
type DetectFraudRequest struct {
	ClaimId string `json:"claim_id"`
	PolicyId string `json:"policy_id"`
	ClaimData map[string]interface{} `json:"claim_data,omitempty"`
}
