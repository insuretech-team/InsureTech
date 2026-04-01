package models


// HealthConditionsRetrievalResponse represents a health_conditions_retrieval_response
type HealthConditionsRetrievalResponse struct {
	Conditions []*ConditionMultiplier `json:"conditions,omitempty"`
}
