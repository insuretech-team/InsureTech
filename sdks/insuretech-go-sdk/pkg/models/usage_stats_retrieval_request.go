package models


// UsageStatsRetrievalRequest represents a usage_stats_retrieval_request
type UsageStatsRetrievalRequest struct {
	ApiKeyId string `json:"api_key_id"`
	EndDate string `json:"end_date,omitempty"`
	StartDate string `json:"start_date,omitempty"`
}
