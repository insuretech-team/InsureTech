package models


// UsageStatsRetrievalRequest represents a usage_stats_retrieval_request
type UsageStatsRetrievalRequest struct {
	EndDate string `json:"end_date,omitempty"`
	ApiKeyId string `json:"api_key_id"`
	StartDate string `json:"start_date,omitempty"`
}
