package models


// StreamStatsResponse represents a stream_stats_response
type StreamStatsResponse struct {
	Stats *ConnectionStats `json:"stats,omitempty"`
}
