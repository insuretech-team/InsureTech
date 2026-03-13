package models


// RenewalGracePeriodRetrievalResponse represents a renewal_grace_period_retrieval_response
type RenewalGracePeriodRetrievalResponse struct {
	GracePeriod *GracePeriod `json:"grace_period,omitempty"`
	Error *Error `json:"error,omitempty"`
}
