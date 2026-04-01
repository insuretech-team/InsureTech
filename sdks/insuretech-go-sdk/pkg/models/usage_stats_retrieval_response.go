package models


// UsageStatsRetrievalResponse represents a usage_stats_retrieval_response
type UsageStatsRetrievalResponse struct {
	AvgResponseTimeMs float64 `json:"avg_response_time_ms,omitempty"`
	FailedRequests string `json:"failed_requests,omitempty"`
	RequestsByEndpoint map[string]interface{} `json:"requests_by_endpoint,omitempty"`
	RequestsByStatusCode map[string]interface{} `json:"requests_by_status_code,omitempty"`
	SuccessfulRequests string `json:"successful_requests,omitempty"`
	TotalRequests string `json:"total_requests,omitempty"`
}
