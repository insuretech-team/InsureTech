package models


// CheckFraudRequest represents a check_fraud_request
type CheckFraudRequest struct {
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Data map[string]interface{} `json:"data,omitempty"`
}
