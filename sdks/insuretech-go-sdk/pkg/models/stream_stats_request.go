package models


// StreamStatsRequest represents a stream_stats_request
type StreamStatsRequest struct {
	RoomId string `json:"room_id"`
	PeerId string `json:"peer_id"`
	IntervalSeconds int `json:"interval_seconds,omitempty"`
}
