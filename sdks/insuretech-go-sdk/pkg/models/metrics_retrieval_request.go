package models


// MetricsRetrievalRequest represents a metrics_retrieval_request
type MetricsRetrievalRequest struct {
	MetricTypes []*MetricType `json:"metric_types"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
}
