package models


// DetectFraudRequest represents a detect_fraud_request
type DetectFraudRequest struct {
	ClaimData map[string]interface{} `json:"claim_data,omitempty"`
	ClaimId string `json:"claim_id"`
	PolicyId string `json:"policy_id"`
}
