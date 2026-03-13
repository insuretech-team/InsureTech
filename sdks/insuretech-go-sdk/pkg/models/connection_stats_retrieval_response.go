package models


// ConnectionStatsRetrievalResponse represents a connection_stats_retrieval_response
type ConnectionStatsRetrievalResponse struct {
	Stats *ConnectionStats `json:"stats,omitempty"`
}
