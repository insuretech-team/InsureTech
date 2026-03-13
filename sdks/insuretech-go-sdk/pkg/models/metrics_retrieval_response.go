package models


// MetricsRetrievalResponse represents a metrics_retrieval_response
type MetricsRetrievalResponse struct {
	Metrics []*AggregatedMetric `json:"metrics,omitempty"`
	Error *Error `json:"error,omitempty"`
}
