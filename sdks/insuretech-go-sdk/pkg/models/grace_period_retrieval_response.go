package models


// GracePeriodRetrievalResponse represents a grace_period_retrieval_response
type GracePeriodRetrievalResponse struct {
	GracePeriod *GracePeriod `json:"grace_period,omitempty"`
	Error *Error `json:"error,omitempty"`
}
