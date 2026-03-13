package models


// PolicyEvaluationRequest represents a policy_evaluation_request
type PolicyEvaluationRequest struct {
	Action string `json:"action"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	UserId string `json:"user_id"`
	Resource string `json:"resource,omitempty"`
}
