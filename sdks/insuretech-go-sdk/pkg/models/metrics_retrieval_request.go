package models


// MetricsRetrievalRequest represents a metrics_retrieval_request
type MetricsRetrievalRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	MetricTypes []*MetricType `json:"metric_types"`
	StartDate string `json:"start_date,omitempty"`
}
