package models


// CheckFraudRequest represents a check_fraud_request
type CheckFraudRequest struct {
	Data map[string]interface{} `json:"data,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}
